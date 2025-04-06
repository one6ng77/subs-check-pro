package app

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/beck-8/subs-check/app/monitor"
	"github.com/beck-8/subs-check/assets"
	"github.com/beck-8/subs-check/check"
	"github.com/beck-8/subs-check/config"
	"github.com/beck-8/subs-check/save"
	"github.com/beck-8/subs-check/utils"
	"github.com/fsnotify/fsnotify"
)

// App 结构体用于管理应用程序状态
type App struct {
	configPath string
	interval   int
	watcher    *fsnotify.Watcher
}

// New 创建新的应用实例
func New() *App {
	configPath := flag.String("f", "", "配置文件路径")
	flag.Parse()

	return &App{
		configPath: *configPath,
	}
}

// Initialize 初始化应用程序
func (app *App) Initialize() error {
	// 初始化配置文件路径
	if err := app.initConfigPath(); err != nil {
		return fmt.Errorf("初始化配置文件路径失败: %w", err)
	}

	// 加载配置文件
	if err := app.loadConfig(); err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	// 初始化配置文件监听
	if err := app.initConfigWatcher(); err != nil {
		return fmt.Errorf("初始化配置文件监听失败: %w", err)
	}

	app.interval = config.GlobalConfig.CheckInterval

	if config.GlobalConfig.ListenPort != "" {
		if err := app.initHttpServer(); err != nil {
			return fmt.Errorf("初始化HTTP服务器失败: %w", err)
		}
	}

	if config.GlobalConfig.SubStorePort != "" {
		if runtime.GOOS == "linux" && runtime.GOARCH == "386" {
			slog.Warn("node不支持Linux 32位系统，不启动sub-store服务")
		}
		go assets.RunSubStoreService()
		// 求等吗得，日志会按预期顺序输出
		time.Sleep(500 * time.Millisecond)
	}

	// 启动内存监控
	monitor.StartMemoryMonitor()

	return nil
}

// Run 运行应用程序主循环
func (app *App) Run() {
	defer func() {
		app.watcher.Close()
	}()

	slog.Info(fmt.Sprintf("进度展示: %v", config.GlobalConfig.PrintProgress))

	for {
		if err := app.checkProxies(); err != nil {
			slog.Error(fmt.Sprintf("检测代理失败: %v", err))
			os.Exit(1)
		}

		nextCheck := time.Now().Add(time.Duration(app.interval) * time.Minute)
		slog.Info(fmt.Sprintf("下次检查时间: %s", nextCheck.Format("2006-01-02 15:04:05")))
		debug.FreeOSMemory()
		time.Sleep(time.Duration(app.interval) * time.Minute)
	}
}

// checkProxies 执行代理检测
func (app *App) checkProxies() error {
	slog.Info("开始检测代理")

	results, err := check.Check()
	if err != nil {
		return fmt.Errorf("检测代理失败: %w", err)
	}
	// 将成功的节点添加到全局中，暂时内存保存
	if config.GlobalConfig.KeepSuccessProxies {
		for _, result := range results {
			if result.Proxy != nil {
				config.GlobalProxies = append(config.GlobalProxies, result.Proxy)
			}
		}
	}

	slog.Info("检测完成")
	save.SaveConfig(results)
	utils.SendNotify(len(results))
	utils.UpdateSubs()
	return nil
}
