package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/server"
)

func StartDeploy(ctx *gin.Context, s *server.Server) {
	var req schema.Template
	filename := ctx.Param("name")

	if filename != "" {
		if !validateFilename(filename) {
			ctx.JSON(400, gin.H{"error": "无效的文件名"})
			return
		}
		// 从文件中读取数据
		data, err := readFlowDataFromFile(filename + ".json")
		if err != nil {
			ctx.JSON(404, gin.H{"error": err.Error()})
		}
		req = data
	} else {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": "无效的请求参数"})
			return
		}
	}
	log.Println(req)
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
