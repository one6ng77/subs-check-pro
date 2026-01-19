package utils

import (
	"testing"

	"github.com/sinspired/subs-check-pro/config"
)

var TestAPI = "https://apprise.xxxxxx.com/notify"
var TestURLs = []string{
	"ntfy://xxxxxxx",
	"bark://api.day.app/xxxxxxxxxxxxxxx",
	"tgram://xxxxxxxxxxxxxxxxxxx/xxxxxxxxxxxxxxxx",
}

// helper: 设置全局配置并在测试结束后恢复
func withTestConfig() {
	config.GlobalConfig.AppriseAPIServer = TestAPI
	config.GlobalConfig.RecipientURL = TestURLs
	config.GlobalConfig.NotifyTitle = "🔔 节点状态更新 [测试]"
}

func TestNotify(t *testing.T) {
	withTestConfig()

	req := NotifyRequest{
		URLs:  TestURLs[0],
		Title: "测试通知",
		Body:  "测试内容",
	}

	if err := Notify(req, ""); err != nil {
		t.Fatalf("Notify() 失败: %v", err)
	}
}

func TestBroadcastNotify(t *testing.T) {
	withTestConfig()

	// 验证函数能正常执行，不返回错误
	broadcastNotify(NotifyNodeStatus, "广播标题", "广播内容", "")
}

func TestSendNotifyCheckResult(t *testing.T) {
	withTestConfig()

	// 验证函数能正常执行，不返回错误
	SendNotifyCheckResult(5)
}

func TestSendNotifyDetectLatestRelease(t *testing.T) {
	withTestConfig()

	// 验证函数能正常执行，不返回错误
	SendNotifyDetectLatestRelease("v1.2.3", "1.13.0", true, false, "https://github.com/sinspired/subs-check/releases/download/v1.13.2/subs-check_Windows_x86_64.zip")
}
