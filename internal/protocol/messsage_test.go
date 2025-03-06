package protocol

import (
	"logicflow-deploy/internal/schema"
	"testing"
)

func TestMessage(t *testing.T) {
	stat, _ := NewMessage(MsgTaskResult, "1", "start", "startNode", schema.NodeStateSuccess)
	t.Log(stat)
	var a schema.NodeStatus
	_ = UnMarshalPayload(stat.Payload, &a)
	t.Log(a)
}
