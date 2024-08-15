package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
	"trading-ace/src/request"
	"trading-ace/src/response"
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
	var query request.GetRewordHistoryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"exception": err.Error()})
		return
	}

	log.Println("Start Time: ", query.StartTime)
	log.Println("End Time: ", query.EndTime)

	startTime, err := time.Parse(time.RFC3339, query.StartTime)
	endTime, err := time.Parse(time.RFC3339, query.EndTime)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"exception": err.Error()})
		return
	}

	rewardRecords, err := r.rewardService.GetRewardHistory(query.User, startTime, endTime.Sub(startTime))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exception": err.Error()})
		return
	}

	pointHistoryCollection := response.CreatePointHistoryCollection(&rewardRecords)
	c.JSON(http.StatusOK, pointHistoryCollection)
}
