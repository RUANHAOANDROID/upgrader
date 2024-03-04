package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"upgrader/config"
)

func main() {
	fmt.Println("auto update!")
	conf, err := config.Load("config.yml")
	if err != nil {
		config.CreateEmpty().Save("config.yml")
		panic("请完善配置config.yml")
	}
	fmt.Println("Timer interval", conf.Timer)
	// ---------------runner -------------
	ctx, cancel := context.WithCancel(context.Background())
	// 监听系统中断信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	StartTimer(conf, ctx, cancel)
}
