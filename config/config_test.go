package config_test

import (
	"fmt"
	"testing"

	"allandeng.cn/allandeng/go-blog/config"
)

func TestConfig(t *testing.T) {
	err := config.LoadConfig("")
	if err != nil {
		fmt.Printf("load config err: %s", err)
	}
	fmt.Println(config.GlobalConfig)
}
