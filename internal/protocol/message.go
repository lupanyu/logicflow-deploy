package protocol

import (
	"encoding/json"
	"time"
)

type MessageType int

const (
	MsgRegister         MessageType = iota + 1
	MsgRegisterResponse             // 注册响应
	MsgJavaDeploy                   // 部署java应用  server->agent
	MsgWebDeploy                    // 部署web应用   server->agent
	MsgShellDeploy                  // shell脚本来部署应用 server->agent
	MsgTask                         // 部署任务控制  server->agent
	MsgHeartbeat                    // 心跳检测   	 server->agent  ,agent->server
	MsgTaskStep                     // 任务详细步骤以及日志信息  agent->server
	MsgTaskResult                   // 任务结果		agent->server
	MsgError                        // 错误信息		agent->server
)

// 部署消息  msgjavadeploy msgwebdeply msgtask
type Message struct {
	Type            MessageType     `json:"type"`
	FlowExecutionID string          `json:"flowExecutionID"`   // flowExecution ID 实际执行flow时生成的id， 或者为空
	AgentID         string          `json:"agentID,omitempty"` // agentID, 一般是hostname
	NodeID          string          `json:"nodeID,omitempty"`  // flow中的nodeID
	Payload         json.RawMessage `json:"payload,omitempty"`
	Timestamp       int64           `json:"timestamp"`
}

func NewMessage(t MessageType, flowId, agentID, nodeID string, payload interface{}) (Message, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Message{}, err
	}
	return Message{
		Type:            t,
		FlowExecutionID: flowId,
		AgentID:         agentID,
		NodeID:          nodeID,
		Payload:         data,
		Timestamp:       time.Now().Unix(),
	}, nil
}

func (m *Message) UpdatePayload(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	m.Payload = data
	return nil
}

func UnMarshalPayload(payload json.RawMessage, data any) error {
	return json.Unmarshal(payload, data)
}
