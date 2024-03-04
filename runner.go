package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"upgrader/pkg"
)

var (
	IsRunning bool
	cmd       *exec.Cmd
)

func RunScript(ctx context.Context) {
	pkg.Log.Printf("RunScript ---")

	// 如果有脚本正在运行，则先停止
	if IsRunning {
		stopScript()
	}

	// 标记当前脚本正在运行
	IsRunning = true

	// 初始化cmd变量
	cmd = exec.CommandContext(ctx, "./runner/app/bin/ledshowktfw")
	logger := log.New(os.Stdout, "cmd", log.LstdFlags)
	// 将命令的标准输出重定向到日志记录器
	cmd.Stdout = logger.Writer()
	// 开始执行命令
	if err := cmd.Run(); err != nil {
		pkg.Log.Printf("Failed to start script: %v\n", err)
		IsRunning = false
		return
	}

	// 监听取消信号，一旦收到取消信号，就终止脚本执行
	go func() {
		<-ctx.Done()
		pkg.Log.Println("Cancellation signal received, terminating script...")
		if err := stopScript(); err != nil {
			pkg.Log.Printf("Failed to stop script: %v\n", err)
		}
		// 标记脚本执行完成
		IsRunning = false
	}()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		pkg.Log.Printf("Script execution error: %v\n", err)
	} else {
		pkg.Log.Println("Script execution completed successfully")
	}
}

func stopScript() error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	// 向脚本进程发送SIGTERM信号，优雅地关闭
	if err := cmd.Process.Signal(os.Kill); err != nil {
		return fmt.Errorf("failed to send SIGINT to process: %v", err)
	}

	// 等待脚本进程退出
	_, err := cmd.Process.Wait()
	if err != nil {
		return fmt.Errorf("failed to wait for process exit: %v", err)
	}

	return nil
}
