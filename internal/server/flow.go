package server

import (
	"encoding/json"
	"log"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
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

func LoadFlowModel(filename string) (schema.Template, error) {
	// 获取目录下的所有文件
	files, err := os.OpenFile(FlowModelDir+filename, os.O_RDONLY, 0644)
	if err != nil {
		return schema.Template{}, err
	}
	defer files.Close()
	var data schema.Template
	jsonParser := json.NewDecoder(files)
	if err := jsonParser.Decode(&data); err != nil {
		return schema.Template{}, err
	}
	return data, nil
}

// 找到当前节点起，下一批能执行的节点
func NextNodes(flowData schema.Template, currentNodeId string) []schema.Node {
	result := []schema.Node{}
	// 遍历所有边，找到当前节点的下一个节点
	for _, edges := range flowData.Edges {
		if edges.SourceNodeId == currentNodeId {
			result = append(result, getNodeById(flowData, edges.TargetNodeId))
		}
	}
	return result
}

func getNodeById(flowData schema.Template, nodeId string) schema.Node {
	for _, node := range flowData.Nodes {
		if node.ID == nodeId {
			return node
		}
	}
	return schema.Node{}
}

func getStartNode(flowData schema.Template) schema.Node {
	for _, node := range flowData.Nodes {
		if node.Type == "start" {
			return node
		}
	}
	return schema.Node{}
}

func CheckDependency(flowData schema.Template, currentNode schema.Node, f schema.FlowExecution) bool {
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
		log.Printf("[%s] 查到mem中存储的依赖节点状态信息是：%v", utils.GetCallerInfo(), status)
		if status.Status != schema.NodeStateSuccess {
			// 依赖的节点没有执行成功
			log.Printf("%s depends on %s : %s", status.Status, dep.ID, status.Status)
			return false
		}
	}
	return true
}
