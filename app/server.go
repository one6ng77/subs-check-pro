package app

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/beck-8/subs-check/config"
	"github.com/beck-8/subs-check/save/method"
	"github.com/gin-gonic/gin"
)

// initHttpServer 初始化HTTP服务器
func (app *App) initHttpServer() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	saver, err := method.NewLocalSaver()
	if err != nil {
		return fmt.Errorf("获取http监听目录失败: %w", err)
	}
	router.Static("/", saver.OutputPath)
	go func() {
		for {
			if err := router.Run(config.GlobalConfig.ListenPort); err != nil {
				slog.Error(fmt.Sprintf("HTTP服务器启动失败，正在重启中: %v", err))
			}
			time.Sleep(30 * time.Second)
		}
	}()
	slog.Info("HTTP服务器启动", "port", config.GlobalConfig.ListenPort)
	return nil
}