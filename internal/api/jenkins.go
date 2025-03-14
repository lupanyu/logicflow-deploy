package api

import (
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/config"
	"net/http"
)

func GetJenkinsData(ctx *gin.Context) {
	response := map[string][]string{}
	for _, jenkins := range config.JenkinsConnections {
		response[jenkins.Name] = jenkins.Jobs
	}
	ctx.JSON(http.StatusOK, response)
}
