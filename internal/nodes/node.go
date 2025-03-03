package nodes

import (
	"context"
	"logicflow-deploy/internal/schema"
)

type NodeExecutor interface {
	Execute(ctx context.Context, state chan schema.TaskStep) // chan使用值传递方式 更新发送给chan，在外部接收
	NodeType() string
	AgentId() string
}
