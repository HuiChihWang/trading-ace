package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"trading-ace/src/service"
)

type RewardController interface {
	GetRewardHistoryOfUser(c *gin.Context)
}

type rewardController struct {
	rewardService service.RewardService
}

var (
	rewardControllerInstance *rewardController
	rewardControllerOnce     sync.Once
)

func GetRewardControllerInstance() RewardController {
	rewardControllerOnce.Do(func() {
		rewardControllerInstance = &rewardController{
			rewardService: service.NewRewardService(),
		}
	})
	return rewardControllerInstance
}

func (r *rewardController) GetRewardHistoryOfUser(c *gin.Context) {
	userId := c.Param("userId")

	rewardRecords, err := r.rewardService.GetRewardHistory(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exception": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rewardRecords)
}
