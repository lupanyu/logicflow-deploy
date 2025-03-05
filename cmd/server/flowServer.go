package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/routes"
	"logicflow-deploy/internal/schema"
	"logicflow-deploy/internal/server"
)

func HandleFlowSave(c *gin.Context) {
	var flowData schema.FlowData
	if err := c.ShouldBindJSON(&flowData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 转换为原始JSON
	rawData, _ := json.Marshal(flowData)

	// 保存文件
	if err := server.SaveFlowModel(rawData); err != nil {
		c.JSON(500, gin.H{"error": "保存失败"})
		return
	}

	c.JSON(200, gin.H{"message": "保存成功"})
}

func main() {

	s := server.NewServer()
	r := gin.Default()
	routes.RegisterAPIRoutes(r, s)
	s.SetHttp(r)
	s.Start("0.0.0.0", 8080)
}
