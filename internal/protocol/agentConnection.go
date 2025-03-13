package protocol

import (
	"github.com/dromara/carbon/v2"
	"github.com/gorilla/websocket"
)

// 定义 AgentConnection 结构体
type AgentConnection struct {
	Conn       *websocket.Conn
	LastActive carbon.Carbon
	Status     agentStatus
}
