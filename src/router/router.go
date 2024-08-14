package router

import (
	"github.com/gin-gonic/gin"
	"trading-ace/src/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/tasks/:userId", controller.GetTaskControllerInstance().GetTasksOfUser)
		apiRoutes.GET("/rewards/:userId", controller.GetRewardControllerInstance().GetRewardHistoryOfUser)
	}

	return r
}
