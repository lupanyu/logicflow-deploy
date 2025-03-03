package nodes

import (
	"context"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
)

type WebNodeExecuter struct {
	properties schema.WebProperties
	agent      *protocol.AgentConnection
}

func (w *WebNodeExecuter) Execute(ctx context.Context, state chan schema.TaskStep) {

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
