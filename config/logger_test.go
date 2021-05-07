package config_test

import (
	"testing"

	"allandeng.cn/allandeng/go-blog/config"
)

func TestLogger(t *testing.T) {
	config.LoadConfig("")
	config.InitLogger()
	log := config.Logger
	if log == nil {
		t.Errorf("log nil")
		return
	}
	log.Debug("go-logging日志打印测试")
	log.Info("go-logging日志打印测试")

	log.Warning("go-logging日志打印测试")
	log.Error("go-logging日志打印测试")
}
