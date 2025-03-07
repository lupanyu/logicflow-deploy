package api

import (
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/server"
)

func GetStatus(ctx *gin.Context, s *server.Server) {
	flowID := ctx.Param("id")
	flowExecution, ok := s.GetFlowExecution(flowID)
	if !ok {
		ctx.JSON(400, gin.H{"error": "Flow execution not found"})
	}
	ctx.JSON(200, flowExecution)
}
