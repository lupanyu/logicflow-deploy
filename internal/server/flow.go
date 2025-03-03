package server

import (
	"encoding/json"
	"fmt"
	"logicflow-deploy/internal/schema"
	"os"
	"path/filepath"
	"time"
)

const FlowModelDir = "./storage/flow-model/"

func SaveFlowModel(data []byte) error {
	// 确保目录存在
	if err := os.MkdirAll(FlowModelDir, 0755); err != nil {
		return err
	}

	// 生成带时间戳的文件名
	filename := filepath.Join(FlowModelDir, "flow-"+time.Now().Format("20060102-150405")+".json")

	// 写入文件
	return os.WriteFile(filename, data, 0644)
}

func LoadFlowModel(filename string) (schema.FlowData, error) {
	// 获取目录下的所有文件
	files, err := os.OpenFile(FlowModelDir+filename, os.O_RDONLY, 0644)
	if err != nil {
		return schema.FlowData{}, err
	}
	defer files.Close()
	var data schema.FlowData
	jsonParser := json.NewDecoder(files)
	if err := jsonParser.Decode(&data); err != nil {
		return schema.FlowData{}, err
	}
	return data, nil
}

// 找到当前节点起，下一批能执行的节点
func NextNodes(flowData schema.FlowData, currentNode schema.Node) []schema.Node {
	result := []schema.Node{}
	// 遍历所有边，找到当前节点的下一个节点
	for _, edges := range flowData.Edges {
		if edges.SourceNodeId == currentNode.ID {
			result = append(result, getNodeById(flowData, edges.TargetNodeId))
		}
	}
	return result
}

func getNodeById(flowData schema.FlowData, nodeId string) schema.Node {
	for _, node := range flowData.Nodes {
		if node.ID == nodeId {
			return node
		}
	}
	return schema.Node{}
}

func getStartNode(flowData schema.FlowData) schema.Node {
	for _, node := range flowData.Nodes {
		if node.Type == "start" {
			return node
		}
	}
	return schema.Node{}
}

func CheckDependency(flowData schema.FlowData, currentNode schema.Node, f schema.FlowExecution) bool {
	// 遍历所有边，找到当前节点的上一个节点
	dependency := []schema.Node{}
	for _, edges := range flowData.Edges {
		if edges.TargetNodeId == currentNode.ID {
			dep := getNodeById(flowData, edges.SourceNodeId)
			dependency = append(dependency, dep)
		}
	}
	for _, dep := range dependency {
		status := f.NodeResults[dep.ID]
		if status.Status != schema.NodeStateSuccess {
			// 依赖的节点没有执行成功
			fmt.Printf("%s depends on %s : %s", status.Status, dep.ID, status.Error)
			return false
		}
	}
	return true
}
