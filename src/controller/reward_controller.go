package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trading-ace/src/service"
)

func GetRewardHistoryOfUser(c *gin.Context) {
	userId := c.Param("userId")

	rewardService := service.NewRewardService()
	rewardRecords, err := rewardService.GetRewardHistory(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rewardRecords)
}
