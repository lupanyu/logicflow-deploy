package nodes

import (
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
)

type WebNodeExecuter struct {
	properties schema.WebProperties
	agent      *protocol.AgentConnection
}

func (w *WebNodeExecuter) Execute(flowExecutionID, nodeID string, ch chan schema.TaskStep, result chan protocol.Message) {
	stat := schema.TaskStep{
		FlowExecutionID: flowExecutionID,
		NodeID:          nodeID,
		Status:          schema.TaskStateSuccess,
		Setup:           "发送web部署指令",
		AgentID:         w.properties.Host,
		Output:          schema.NewOutLog(schema.LevelInfo, "发送web部署指令", "..."),
	}

	// 执行部署命令
	err := w.deploy()
	if err != nil {
		stat.Status = schema.TaskStateFailed
		stat.Error = schema.NewOutLog(schema.LevelError, "发送web部署指令", err.Error())
		log.Printf("[%s] 向%s发送部署指令异常，参数是：%v， 错误是: %v", utils.GetCallerInfo(), w.properties.Host, w.properties, err.Error())
	} else {
		log.Printf("[%s] 向%s发送部署指令成功，参数是：%v", utils.GetCallerInfo(), w.properties.Host, w.properties)
	}
	ch <- stat
}

func NewWebNodeExecuter(data schema.WebProperties, agent *protocol.AgentConnection) *WebNodeExecuter {

	return &WebNodeExecuter{
		properties: data,
		agent:      agent,
	}
}

// 向agent发送部署命令
func (w *WebNodeExecuter) deploy() error {
	return w.agent.WriteJSON(w.properties)
}
