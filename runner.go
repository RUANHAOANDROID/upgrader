package main

import (
	"context"
	"os/exec"
	"upgrader/pkg"
)

func runScript(ctx context.Context) {
	// 在这里替换为你要执行的脚本和参数
	cmd := exec.CommandContext(ctx, "./runner/app/bin/ledshowktfw")

	// 开始执行命令
	if err := cmd.Start(); err != nil {
		pkg.Log.Printf("Failed to start script: %v\n", err)
		return
	}

	// 监听取消信号，一旦收到取消信号，就终止脚本执行
	go func() {
		<-ctx.Done()
		pkg.Log.Println("Cancellation signal received, terminating script...")
		if err := cmd.Process.Kill(); err != nil {
			pkg.Log.Printf("Failed to kill script: %v\n", err)
		}
	}()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		pkg.Log.Printf("Script execution error: %v\n", err)
	} else {
		pkg.Log.Println("Script execution completed successfully")
	}
}
