package api

import (
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/server"
)

func StartDeploy(ctx *gin.Context, s *server.Server) {
	var req schema.FlowData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}
	process, err := server.NewFlowProcessor(req, s)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	flowExecution := process.ExecuteFlow(s)
	ctx.JSON(200, flowExecution)
}

func CancelDeploy(ctx *gin.Context) {

}
