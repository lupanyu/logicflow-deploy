package protocol

import (
	"github.com/gorilla/websocket"
	"time"
)

// 定义 AgentConnection 结构体
type AgentConnection struct {
	Conn       *websocket.Conn
	LastActive time.Time
	Status     agentStatus
}
