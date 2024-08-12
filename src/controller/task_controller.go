package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trading-ace/src/service"
)

var taskService service.TaskService = service.NewTaskService()

func GetTasksOfUser(c *gin.Context) {
	userId := c.Param("userId")
	tasks, err := taskService.GetTasksOfUser(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exception": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
