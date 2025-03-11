package server

import (
	"logicflow-deploy/internal/schema"
)

// 在Server结构体中添加方法
func (s *Server) GetFlowExecution(flowID string) (schema.FlowExecution, bool) {
	// 从存储中获取流程执行状态
	execution, ok := s.stateStorage.Get(flowID)
	if !ok {
		return schema.FlowExecution{}, false
	}
	return execution, true
}

func (s *Server) GetAllFlowExecution() []schema.FlowExecution {
	return s.stateStorage.GetAll()
}
