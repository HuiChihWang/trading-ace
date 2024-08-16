package router

import (
	"github.com/gin-gonic/gin"
	"trading-ace/src/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/tasks", controller.GetTaskControllerInstance().SearchTasks)
		apiRoutes.GET("/reward-history", controller.GetRewardControllerInstance().GetRewardHistoryOfUser)
	}

	return r
}
