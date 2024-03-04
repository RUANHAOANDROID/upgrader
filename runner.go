package main

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"upgrader/pkg"
)

var (
	IsRunning bool
	cmd       *exec.Cmd
)

func RunScript(ctx context.Context) {
	Kill6688()
	// 启动ledshowktfw服务
	cmd := exec.Command("./runner/app/bin/ledshowktfw")

	// 创建管道来捕获命令的输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		pkg.Log.Println("Error:", err)
		IsRunning = false
		return
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		pkg.Log.Println("Error:", err)
		IsRunning = false
		return
	}
	IsRunning = true
	// 读取并输出命令的输出
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		pkg.Log.Println(scanner.Text())
	}

	// 检查是否发生错误
	if err := scanner.Err(); err != nil {
		pkg.Log.Println("Error:", err)
	}

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		pkg.Log.Println("Error:", err)
	}
}

func Kill6688() {
	// 执行Shell脚本
	cmd := exec.Command("/bin/sh", "-c", "./kill6688.sh")
	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 输出执行结果
	fmt.Println(string(output))
	IsRunning = false
}
