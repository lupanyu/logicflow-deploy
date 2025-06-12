package main

import (
	"flag"
	"log"
	"logicflow-deploy/internal/agent"
	"logicflow-deploy/internal/utils"
	"time"
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
	for {
		// 初始化 agent
		da := agent.NewDeploymentAgent(serverURL)
		// 建立 WebSocket 连接
		_ = da.Connect()
		retry := 0
		select {
		case <-da.Done:
			// 清理资源
			da.Cleanup()
			log.Println("连接已关闭，尝试重新连接...")
			// 指数退避策略
			waitTime := time.Duration(retry*retry) * time.Second
			if waitTime > 30*time.Second {
				waitTime = 30 * time.Second
			}
			log.Printf(" [%s]等待%.0f秒后尝试第%d次重连...", utils.GetCallerInfo(), waitTime.Seconds(), retry+1)
			time.Sleep(waitTime)
			retry++
			continue
		}
	}
}
