package server

import (
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/schema"
	"net/http"
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

// 在文件底部添加HTTP路由处理函数（需要确保server包有HTTP路由设置）
func (s *Server) handleGetFlowStatus(c *gin.Context) {
	flowID := c.Param("id")
	execution, ok := s.GetFlowExecution(flowID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flow execution not found"})
		return
	}
	c.JSON(http.StatusOK, execution)
}
