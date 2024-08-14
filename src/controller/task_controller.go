package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"trading-ace/src/service"
)

type TaskController interface {
	GetTasksOfUser(c *gin.Context)
}

type taskController struct {
	taskService service.TaskService
}

var (
	taskControllerInstance *taskController
	taskControllerOnce     sync.Once
)

func GetTaskControllerInstance() TaskController {
	taskControllerOnce.Do(func() {
		taskControllerInstance = &taskController{
			taskService: service.NewTaskService(),
		}
	})
	return taskControllerInstance
}

func (t *taskController) GetTasksOfUser(c *gin.Context) {
	userId := c.Param("userId")
	tasks, err := t.taskService.GetTasksOfUser(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exception": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
