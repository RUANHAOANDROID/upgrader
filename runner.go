package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
		StopScript()
	}

	// 初始化cmd变量
	cmd = exec.CommandContext(ctx, "./runner/app/bin/ledshowktfw")
	logger := log.New(os.Stdout, "cmd", log.LstdFlags)
	// 将命令的标准输出重定向到日志记录器
	cmd.Stdout = logger.Writer()
	// 标记当前脚本正在运行
	IsRunning = true
	// 开始执行命令
	if err := cmd.Start(); err != nil {
		pkg.Log.Error("Failed to start script: %v\n", err)
		IsRunning = false
		return
	}

	// 监听取消信号，一旦收到取消信号，就终止脚本执行
	go func() {
		<-ctx.Done()
		pkg.Log.Println("Cancellation signal received, terminating script...")
		if err := StopScript(); err != nil {
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

func StopScript() error {
	kill6688()
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
func kill6688() {
	// 获取监听在端口6688的进程的PID
	pid, err := findProcessID(6688)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 打开进程
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 向进程发送SIGKILL信号
	err = process.Kill()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Process with PID", pid, "has been killed")
}

func findProcessID(port int) (int, error) {
	// 获取所有TCP监听的端口
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0, err
	}

	// 遍历所有监听的端口，找到占用指定端口的进程的PID
	for _, addr := range addrs {
		if addr.Network() == "tcp" {
			tcpAddr := addr.(*net.IPNet)
			if tcpAddr.IP.IsLoopback() {
				continue
			}
			addrString := tcpAddr.IP.String() + ":" + fmt.Sprintf("%d", port)
			conn, err := net.Dial("tcp", addrString)
			if err != nil {
				continue
			}
			conn.Close()

			cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-t")
			output, err := cmd.Output()
			if err != nil {
				return 0, err
			}

			pid := string(output)
			return strconv.Atoi(strings.TrimSpace(pid))
		}
	}
	return 0, fmt.Errorf("Process not found")
}
