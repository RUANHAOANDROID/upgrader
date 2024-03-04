package main

import (
	"time"
	"upgrader/config"
)

func StartTimer(conf *config.Config) {
	timer := time.NewTimer(conf.Timer)
	Auto()
	defer timer.Stop()
	for {
		timer.Reset(conf.Timer)
		select {
		case <-timer.C:
			Auto()
		}
	}
}
