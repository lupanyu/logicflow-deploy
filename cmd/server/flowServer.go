package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/server"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 添加WebSocket升级器配置
var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求，生产环境应限制
	},
}

// HandleFlowExecution 新增WebSocket处理函数
func HandleFlowExecution(s *server.Server, c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "WebSocket升级失败"})
		return
	}
	// defer conn.Close()

	// 设置读取数据超时时间
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	go server.HandleAgentConnection(s, conn)
}
func HandleFlowSave(c *gin.Context) {
	var flowData schema.FlowData
	if err := c.ShouldBindJSON(&flowData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 转换为原始JSON
	rawData, _ := json.Marshal(flowData)

	// 保存文件
	if err := server.SaveFlowModel(rawData); err != nil {
		c.JSON(500, gin.H{"error": "保存失败"})
		return
	}

	c.JSON(200, gin.H{"message": "保存成功"})
}

func main() {
	s := server.NewServer()
	r := gin.Default()
	r.Static("/static", "./static")
	r.POST("/save", HandleFlowSave)

	r.GET("/ws", func(c *gin.Context) { HandleFlowExecution(s, c) })
	s.SetHttp(r)
	s.Start(8080)
}
