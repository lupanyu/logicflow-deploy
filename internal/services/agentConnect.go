package services

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"logicflow-deploy/internal/server"
	"net/http"
	"time"
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

	// 设置读取数据超时时间
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	go server.HandleAgentConnection(s, conn)
}
