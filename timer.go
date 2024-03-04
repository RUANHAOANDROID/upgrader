package main

import (
	"context"
	"time"
	"upgrader/config"
)

func StartTimer(conf *config.Config, ctx context.Context, cancel context.CancelFunc) {
	timer := time.NewTimer(conf.Timer)
	Auto(ctx, cancel)
	defer timer.Stop()
	for {
		timer.Reset(conf.Timer)
		select {
		case <-timer.C:
			Auto(ctx, cancel)
		}
	}
}
