package schema

import (
	"encoding/json"
	"fmt"
	"github.com/dromara/carbon/v2"
)

type Template struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Env         string `json:"env"`
	Nodes       []Node `json:"nodes"`
	Edges       []Edge `json:"edges"`
}

type FlowData struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	X          float64         `json:"x"`
	Y          float64         `json:"y"`
	Properties json.RawMessage `json:"properties"`
	Text       Text            `json:"text,omitempty"`
}

// 反序列化properties
func (n *Node) DeserializationProperties(data interface{}) error {
	return json.Unmarshal(n.Properties, &data)
}

type Edge struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Properties   EdgeProperty `json:"properties"`
	SourceNodeId string       `json:"sourceNodeId"`
	TargetNodeId string       `json:"targetNodeId"`
	StartPoint   Point        `json:"startPoint"`
	EndPoint     Point        `json:"endPoint"`
	PointsList   []Point      `json:"pointsList"`
}

type EdgeProperty struct {
	Color string `json:"color"`
}

type Text struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Value string  `json:"value"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// 每个节点的执行状态，用以判断节点是否在运行、成功、失败等状态,在node 结束后传递这个结构给server
type NodeState struct {
	ID        string         `json:"id"` // flowData node 节点的id
	Type      string         `json:"type"`
	Status    NodeStatus     `json:"status"` // 状态: pending/running/success/failed
	StartTime *carbon.Carbon `json:"startTime"`
	EndTime   *carbon.Carbon `json:"endTime"`
	Logs      string         `json:"logs"`  // 执行日志
	Error     string         `json:"error"` // 错误信息
}

func (n *NodeState) AppendTaskStep(taskStep TaskStep) {
	if n.StartTime == nil {
		now := carbon.Now()
		n.StartTime = &now
	}
	n.Logs += fmt.Sprintf("%s %s \n", taskStep.Output, taskStep.Status)
	if taskStep.Error != "" {
		n.Error += taskStep.Error
	}
}

// FlowExecution 表示整个flow流程，包含所有节点的执行结果
type FlowExecution struct {
	FlowID       string               `json:"flowId"`      // 流程唯一ID
	GlobalStatus NodeStatus           `json:"status"`      // 全局执行状态（新增）
	StartTime    *carbon.Carbon       `json:"startTime"`   // 改为指针类型，可空
	EndTime      *carbon.Carbon       `json:"endTime"`     // 改为指针类型，可空
	Duration     float64              `json:"duration"`    // 新增执行耗时（秒）
	NodeResults  map[string]NodeState `json:"nodeResults"` // key 是 node 的 id，value 是 node 的执行结果
	FlowData     FlowData             `json:"flowData"`    // 原始数据
}

// 持续时间计算方法
func (f *FlowExecution) CalculateDuration() {
	if f.StartTime != nil && f.EndTime != nil {
		f.Duration = float64(int(f.EndTime.DiffAbsInDuration(*f.StartTime)))
	}
}

// 日志条目结构
type LogEntry struct {
	Timestamp carbon.Carbon `json:"timestamp"`
	Message   string        `json:"message"`
	Level     string        `json:"level"` // info/warn/error
}

// 新增错误详情结构
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Stack   string `json:"stack,omitempty"`
}

// 判断流程是否完成
func (f *FlowExecution) IsCompleted() bool {
	return f.GlobalStatus == NodeStateSuccess ||
		f.GlobalStatus == NodeStateFailed ||
		f.GlobalStatus == NodeStateTimeout
}
