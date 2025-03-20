package protocol

import (
	"github.com/dromara/carbon/v2"
	"github.com/gorilla/websocket"
	"sync"
)

// 定义 AgentConnection 结构体
type AgentConnection struct {
	Conn       *websocket.Conn
	LastActive carbon.Carbon
	Status     agentStatus
	lock       sync.RWMutex
}

func (a *AgentConnection) WriteJSON(msg interface{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.Conn.WriteJSON(msg)
}
