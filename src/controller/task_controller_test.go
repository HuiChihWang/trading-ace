package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"trading-ace/mock/service"
	"trading-ace/src/model"
)

type taskControllerTestSuite struct {
	taskController    TaskController
	mockedTaskService *service.MockTaskService
}

func (s *taskControllerTestSuite) setUp(t *testing.T) {
	s.mockedTaskService = service.NewMockTaskService(t)
	s.taskController = &taskController{
		taskService: s.mockedTaskService,
	}
}

func TestGetTasksOfUser(t *testing.T) {
	testSuite := &taskControllerTestSuite{}
	//ID          int          `json:"id"`
	//Status      TaskStatus   `json:"status"`
	//Type        TaskType     `json:"type"`
	//UserID      string       `json:"user_id"`
	//SwapAmount  float64      `json:"swap_amount"`
	//CreatedAt   time.Time    `json:"created_at"`
	//CompletedAt sql.NullTime `json:"completed_at"`
	tasks := []*model.Task{
		{
			ID:         1,
			Status:     model.TaskStatusPending,
			Type:       model.TaskTypeSharedPool,
			UserID:     "test_user_id",
			SwapAmount: 10.0,
		},
		{
			ID:         2,
			Status:     model.TaskStatusDone,
			Type:       model.TaskTypeOnboarding,
			UserID:     "test_user_id",
			SwapAmount: 10000.0,
		},
		{
			ID:         3,
			Status:     model.TaskStatusPending,
			Type:       model.TaskTypeSharedPool,
			UserID:     "test_user_id",
			SwapAmount: 10.0,
		},
	}

	t.Run("GetTasksOfUser", func(t *testing.T) {
		testSuite.setUp(t)

		testSuite.mockedTaskService.EXPECT().
			GetTasksOfUser("test_user_id").
			Return(tasks, nil)

		testResponseWriter := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(testResponseWriter)
		testContext.Params = gin.Params{
			{Key: "userId", Value: "test_user_id"},
		}

		testSuite.taskController.GetTasksOfUser(testContext)

		assert.Equal(t, http.StatusOK, testContext.Writer.Status())

		var tasksFromRes []*model.Task
		err := json.Unmarshal(testResponseWriter.Body.Bytes(), &tasksFromRes)

		assert.Nil(t, err)
		assert.Equal(t, tasks, tasksFromRes)
	})

	t.Run("GetTasksOfUser with error", func(t *testing.T) {
		testSuite.setUp(t)

		testSuite.mockedTaskService.EXPECT().
			GetTasksOfUser("test_user_id").
			Return(nil, assert.AnError)

		testResponseWriter := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(testResponseWriter)
		testContext.Params = gin.Params{
			{Key: "userId", Value: "test_user_id"},
		}

		testSuite.taskController.GetTasksOfUser(testContext)

		assert.Equal(t, http.StatusInternalServerError, testContext.Writer.Status())

		var exception map[string]string
		err := json.Unmarshal(testResponseWriter.Body.Bytes(), &exception)

		assert.Nil(t, err)
		assert.Equal(t, assert.AnError.Error(), exception["exception"])
	})
}
