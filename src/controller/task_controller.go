package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"trading-ace/src/repository"
	"trading-ace/src/request"
	"trading-ace/src/response"
	"trading-ace/src/service"
)

type TaskController interface {
	SearchTasks(c *gin.Context)
}

type taskController struct {
	taskService   service.TaskService
	rewardService service.RewardService
}

var (
	taskControllerInstance *taskController
	taskControllerOnce     sync.Once
)

func GetTaskControllerInstance() TaskController {
	taskControllerOnce.Do(func() {
		taskControllerInstance = &taskController{
			taskService:   service.NewTaskService(),
			rewardService: service.NewRewardService(),
		}
	})
	return taskControllerInstance
}

func (t *taskController) SearchTasks(c *gin.Context) {
	var query request.GetRewordHistoryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"exception": err.Error()})
		return
	}

	tasks, err := t.taskService.SearchTasks(&repository.SearchTasksCondition{
		UserID: query.User,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exception": err.Error()})
		return
	}

	var tasksRes response.TaskCollection
	for _, task := range *tasks {
		distributedPoint := 0.0
		rewardRecord, _ := t.rewardService.GetRewardHistoryByTaskID(task.ID)

		if rewardRecord != nil {
			distributedPoint = rewardRecord.Points
		}

		taskRes := response.NewTask(task, distributedPoint)
		tasksRes = append(tasksRes, taskRes)
	}

	if tasksRes == nil {
		tasksRes = response.TaskCollection{}
	}

	c.JSON(http.StatusOK, &tasksRes)
}
