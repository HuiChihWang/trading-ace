package router

import (
	"github.com/gin-gonic/gin"
	"trading-ace/src/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/rewards/:userId", controller.GetRewardHistoryOfUser)
	}

	return r
}
