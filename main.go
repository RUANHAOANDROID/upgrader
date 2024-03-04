package main

import (
	"fmt"
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
	StartTimer(conf)
}
