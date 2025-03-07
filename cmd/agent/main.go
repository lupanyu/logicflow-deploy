package main

import (
	"flag"
	"log"
	"logicflow-deploy/internal/agent"
	"logicflow-deploy/internal/utils"
)

func main() {
	utils.InitLog()
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
	da.Connect()

}
