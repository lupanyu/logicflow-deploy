package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"logicflow-deploy/internal/config"
	"logicflow-deploy/internal/nodes"
	"logicflow-deploy/internal/protocol"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/utils"
	"time"
)

// FlowProcessor 是流程处理器的数据结构
type FlowProcessor struct {
	FlowID         string
	flowData       schema.FlowData
	executors      map[string]nodes.NodeExecutor
	taskStepChan   chan schema.TaskStep  // 把节点里的每个步骤的状态发送到这里
	taskResultChan chan protocol.Message // 把每个node节点最终的状态发送到这里 payload是NodeStatus
	ctx            context.Context
	cancel         context.CancelFunc
}

func (fp *FlowProcessor) Cancel() {
	log.Printf(" [%s]终止执行flow: %s", utils.GetCallerInfo(), fp.FlowID)
	// 终止所有的executor,TODO
	// 终止所有的taskStepChan
	close(fp.taskStepChan)
	// 更新最终状态为终止 这个在外层做 TODO
	// 终止所有的taskResultChan
	close(fp.taskResultChan)

	fp.cancel()
}

// 初始化处理器
func NewFlowProcessor(flow schema.FlowData, s *Server) (*FlowProcessor, error) {
	ctx, cancel := context.WithCancel(context.Background())
	fp := &FlowProcessor{
		FlowID:         generateFlowID(),
		flowData:       flow,
		executors:      make(map[string]nodes.NodeExecutor),
		taskStepChan:   make(chan schema.TaskStep),
		taskResultChan: make(chan protocol.Message),
		ctx:            ctx,
		cancel:         cancel,
	}
	log.Printf("[%s] 初始化流程处理器: %v", utils.GetCallerInfo(), fp)
	for _, node := range flow.Nodes {
		switch node.Type {
		case "start":
			fp.RegisterExecutor(node.ID, nodes.NewStartNodeExecutor(node))
		case "stop":
			fp.RegisterExecutor(node.ID, nodes.NewEndNodeExecutor(node))
		case "java":
			var data schema.JavaProperties
			// 从node.Properties中获取属性
			err := node.DeserializationProperties(&data)
			if err != nil {
				return nil, err
			}
			log.Println("反序列化java properties", data)

			host, ok := s.GetAgentConnection(data.Host)
			if !ok {
				return nil, fmt.Errorf("java节点%v未找到AgentID: %s", data, data.Host)
			}
			if !s.HandleAgentStatus(data.Host) {
				return nil, fmt.Errorf("AgentID: %s 状态异常", data.Host)
			}
			fp.RegisterExecutor(node.ID, nodes.NewJavaNodeExecutor(data, host))
		case "build":
			fp.RegisterExecutor(node.ID, nodes.NewBuildNodeExecutor(node))
		case "web":
			// 反序列化properties
			var data schema.WebProperties
			// 从node.Properties中获取属性

			err := node.DeserializationProperties(&data)
			if err != nil {
				return nil, err
			}
			log.Println("反序列化web properties", data)

			host, ok := s.GetAgentConnection(data.Host)
			if !ok {
				return nil, fmt.Errorf("web节点%v未找到AgentID: %s", data, data.Host)
			}
			fp.RegisterExecutor(node.ID, nodes.NewWebNodeExecuter(data, host))
		case "jenkins":
			// 反序列化properties
			var data schema.JenkinsProperties
			err := node.DeserializationProperties(&data)
			if err != nil {
				return nil, err
			}
			log.Println("反序列化jenkins properties", data)
			host, ok := config.JenkinsConnections[data.NodeName]
			if !ok {
				return nil, fmt.Errorf("jenkins节点%v未找到AgentID: %s", data, data.NodeName)
			}
			fp.RegisterExecutor(node.ID, nodes.NewJenkinsNodeExecutor(data, host))
		case "shell":
			// 反序列化properties
			var data schema.ShellProperties
			err := node.DeserializationProperties(&data)
			if err != nil {
				return nil, err
			}
			log.Println("反序列化shell properties", data)
			host, ok := s.GetAgentConnection(data.Host)
			if !ok {
				return nil, fmt.Errorf("shell节点%v未找到AgentID: %s", data, data.Host)
			}
			fp.RegisterExecutor(node.ID, nodes.NewShellNodeExecutor(data, host))
		case "end":
			fp.RegisterExecutor(node.ID, nodes.NewEndNodeExecutor(node))
		default:
			log.Printf(" [%s]未注册的节点类型: %s", utils.GetCallerInfo(), node.Type)
			return nil, fmt.Errorf("未注册的节点类型: %s", node.Type)
		}
	}
	// 把fp添加到server的fpMap中
	s.fpMap[fp.FlowID] = fp
	// 初始化存储数据
	flowExecution := schema.FlowExecution{
		FlowID:       fp.FlowID,
		GlobalStatus: schema.NodeStatePending,
		FlowData:     fp.flowData,
		NodeResults:  make(map[string]schema.NodeState),
	}
	// 写入初始node数据
	for _, node := range fp.flowData.Nodes {
		flowExecution.NodeResults[node.ID] = schema.NodeState{
			Status: schema.NodeStatePending,
			Logs:   "",
			Error:  "",
			ID:     node.ID,
			Type:   node.Type,
		}
	}

	s.stateStorage.Save(flowExecution)
	return fp, nil
}

// 生成唯一的FlowID
func generateFlowID() string {
	return uuid.New().String()
}

// // 接收statusChan中的数据 并将数据存到Storage中
func (fp *FlowProcessor) statusFactory(mem Storage, s *Server) {

	for {
		// 如果存储中没有flowExecution 就创建一个,每次更新的数据都从存储中获取
		flowExecution, ok := mem.Get(fp.FlowID)
		if !ok {
			flowExecution = schema.FlowExecution{
				FlowID:      fp.FlowID,
				NodeResults: make(map[string]schema.NodeState),
			}
		}
		log.Println("waiting task status info...")
		select {
		// 收到状态更新
		case taskStep := <-fp.taskStepChan:
			log.Printf("收到 taskstep: %v", taskStep)
			// 把每一步的日志更新到 flowExecution 中
			taskStepData := flowExecution.NodeResults[taskStep.NodeID]
			taskStepData.AppendTaskStep(taskStep)
			flowExecution.NodeResults[taskStep.NodeID] = taskStepData
			mem.Save(flowExecution)
			// 收到节点状态更新
		case state := <-fp.taskResultChan:
			log.Printf("收到 taskResult: %v", state.Payload)
			// 提取关键状态参数信息
			var nodeStatus schema.NodeStatus
			err := protocol.UnMarshalPayload(state.Payload, &nodeStatus)
			if err != nil {
				log.Printf(" [%s]反序列化状态消息失败: %v", utils.GetCallerInfo(), err)
				return
			}
			log.Printf(" [%s]状态工厂收到node节点%s 最后状态是: %v", utils.GetCallerInfo(), state.NodeID, nodeStatus)
			// 更新node节点的状态
			nodeState := flowExecution.NodeResults[state.NodeID]
			nodeState.Status = nodeStatus
			now := time.Now()
			nodeState.EndTime = &now
			// 节点执行成功，更新状态
			flowExecution.NodeResults[state.NodeID] = nodeState
			mem.Save(flowExecution)
			if nodeStatus == schema.NodeStateFailed {
				flowExecution.GlobalStatus = schema.NodeStateFailed
				flowExecution.EndTime = &now
				log.Printf(" [%s]flow: %s 执行失败,结束", utils.GetCallerInfo(), flowExecution.FlowID)
				mem.Save(flowExecution)
				break
			}
			// 如果当前节点是结果是成功的，触发下一个节点的执行
			if nodeStatus == schema.NodeStateSuccess {
				nextNodes := NextNodes(flowExecution.FlowData, state.NodeID)
				if len(nextNodes) == 0 {
					log.Printf(" [%s]flow: %s 没有剩余要执行的节点,结束", utils.GetCallerInfo(), flowExecution.FlowID)
					flowExecution.GlobalStatus = schema.NodeStateSuccess
					flowExecution.EndTime = &now
					mem.Save(flowExecution)
					break
				}
				log.Printf(" [%s]节点%s执行成功，下一个节点是: %v", utils.GetCallerInfo(), state.NodeID, nextNodes)
				for _, nextNode := range nextNodes {
					log.Printf(" [%s]自动执行下一节点: %v", utils.GetCallerInfo(), nextNode)
					go fp.executeNode(nextNode, s)
				}
			} else {
				log.Printf(" [%s]节点%s执行失败，停止执行flowID:%s", utils.GetCallerInfo(), state.NodeID, state.FlowExecutionID)
				break
			}

		}
	}
}

// 执行flow
func (fp *FlowProcessor) ExecuteFlow(server *Server) schema.FlowExecution {
	execution, ok := server.stateStorage.Get(fp.FlowID)
	if !ok {
		return schema.FlowExecution{}
	}
	// 保存flowExecution
	execution.GlobalStatus = schema.NodeStateRunning
	now := time.Now()
	execution.StartTime = &now
	server.stateStorage.Save(execution)

	// 启动状态工厂
	go fp.statusFactory(server.stateStorage, server)
	// 查找开始节点
	var startNode schema.Node
	for _, node := range fp.flowData.Nodes {
		if node.Type == "start" {
			startNode = node
			break
		}
	}
	log.Printf("[%s] 执行节点:%v", utils.GetCallerInfo(), startNode.ID)

	// 执行节点
	go fp.executeNode(startNode, server)
	return execution
}

// 执行单个节点
func (fp *FlowProcessor) executeNode(node schema.Node, server *Server) {

	// 查看存储中是否有当前节点的执行日志
	flowExecution, ok := server.stateStorage.Get(fp.FlowID)
	if !ok {
		log.Printf(" [%s]未找到FlowExecution: %s,停止执行", utils.GetCallerInfo(), fp.FlowID)
		return
	}
	log.Println(flowExecution)
	if CheckDependency(fp.flowData, node, flowExecution) {
		log.Println("依赖节点执行成功")
	} else {
		log.Println("依赖节点未执行成功，停止执行当前节点", node.ID)
		return
	}

	executor, exists := fp.executors[node.ID]
	if !exists {
		log.Printf(" [%s]未找到节点执行器: %s", utils.GetCallerInfo(), node.ID)
		return
	}
	log.Printf("[%s]执行节点:%s", utils.GetCallerInfo(), node.ID)
	executor.Execute(fp.FlowID, node.ID, fp.taskStepChan, fp.taskResultChan)

}

// 注册节点处理器
func (fp *FlowProcessor) RegisterExecutor(nodeId string, executor nodes.NodeExecutor) {
	fp.executors[nodeId] = executor
}
