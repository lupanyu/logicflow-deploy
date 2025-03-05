package nodes

import (
	"log"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
)

type WebNodeExecuter struct {
	properties schema.WebProperties
	agent      *protocol.AgentConnection
}

func (w *WebNodeExecuter) Execute() schema.TaskStep {
	stat := schema.TaskStep{
		Status:  schema.TaskStateSuccess,
		Setup:   "发送部署指令",
		AgentID: w.properties.Host,
		Output:  schema.NewOutLog(schema.LevelInfo, "开始应用部署"),
	}

	// 执行部署命令
	err := w.deploy()
	if err != nil {
		stat.Status = schema.TaskStateFailed
		stat.Error = schema.NewOutLog(schema.LevelError, err.Error())
		log.Println("向%s发送部署指令异常，参数是：%s， 错误是: %v", w.properties.Host, w.properties, err.Error())
	} else {
		log.Println("向%s发送部署指令成功，参数是：%s", w.properties.Host, w.properties)
	}
	return stat
}
func (w *WebNodeExecuter) NodeType() string {
	return "web"
}
func (w *WebNodeExecuter) AgentId() string {
	return w.properties.Host
}
func NewWebNodeExecuter(data schema.WebProperties, agent *protocol.AgentConnection) *WebNodeExecuter {

	return &WebNodeExecuter{
		properties: data,
		agent:      agent,
	}
}

// 向agent发送部署命令
func (w *WebNodeExecuter) deploy() error {
	return w.agent.Conn.WriteJSON(w.properties)
}
