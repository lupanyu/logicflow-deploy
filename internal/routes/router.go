package routes

import (
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/api"
	"logicflow-deploy/internal/server"
	"logicflow-deploy/internal/services"
)

func RegisterAPIRoutes(r *gin.Engine, s *server.Server) {
	// jenkins 任务
	r.GET("/api/v1/jenkins", api.GetJenkinsData)
	// 部署相关路由组
	deployGroup := r.Group("/api/v1/deploy")
	{
		deployGroup.POST("/", func(context *gin.Context) {
			api.StartDeploy(context, s)
		}) // 开始部署
		deployGroup.POST("/:name", func(context *gin.Context) {
			api.StartDeploy(context, s)
		}) // 开始部署
		deployGroup.GET("/:id", func(context *gin.Context) {
			api.GetStatus(context, s)
		}) // 获取部署的状态
		deployGroup.GET("/", func(context *gin.Context) {
			api.GetDeployHistoryList(context, s)
		})
	}
	flowGroup := r.Group("/api/v1/flow")
	{
		flowGroup.POST("/", api.CreateFlowData)
		flowGroup.GET("/:name", api.GetFlow)
		flowGroup.PUT("/:name", api.UpdateFlow)
		flowGroup.DELETE("/:name", api.DeleteFlow)
		flowGroup.GET("/", api.ListFlow)
	}

	r.GET("/ws", func(c *gin.Context) { services.HandleFlowExecution(s, c) })
}
