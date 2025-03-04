package main

import (
	"flag"
	"log"
	"logicflow-deploy/internal/agent"
)

func main() {
	// 解析命令行参数
	var serverURL string
	flag.StringVar(&serverURL, "server", "", "LogicFlow server URL (ws:// or wss://)")
	flag.Parse()

	if serverURL == "" {
		log.Fatal("必须通过 --server 参数指定服务器地址")
	}

	// 初始化 agent
	da := agent.NewDeploymentAgent(serverURL)

	// 建立 WebSocket 连接
	if err := da.Connect(); err != nil {
		log.Fatalf("连接服务器失败: %v", err)
	}

}
