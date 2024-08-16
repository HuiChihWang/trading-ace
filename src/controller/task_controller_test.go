package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"trading-ace/mock/service"
	"trading-ace/src/model"
	"trading-ace/src/repository"
	"trading-ace/src/response"
)

type taskControllerTestSuite struct {
	taskController      TaskController
	mockedTaskService   *service.MockTaskService
	mockedRewardService *service.MockRewardService
}

func (s *taskControllerTestSuite) setUp(t *testing.T) {
	s.mockedTaskService = service.NewMockTaskService(t)
	s.mockedRewardService = service.NewMockRewardService(t)
	s.taskController = &taskController{
		taskService:   s.mockedTaskService,
		rewardService: s.mockedRewardService,
	}
}

func TestGetTasksOfUser(t *testing.T) {
	testSuite := &taskControllerTestSuite{}
	createTime, _ := time.Parse(time.RFC3339, "2021-01-01T05:00:00Z")
	tasks := []*model.Task{
		{
			ID:         1,
			Status:     model.TaskStatusPending,
			Type:       model.TaskTypeSharedPool,
			UserID:     "test_user_id",
			SwapAmount: 10.0,
			CreatedAt:  createTime,
		},
		{
			ID:         2,
			Status:     model.TaskStatusDone,
			Type:       model.TaskTypeOnboarding,
			UserID:     "test_user_id",
			SwapAmount: 10000.0,
			CreatedAt:  createTime,
		},
		{
			ID:         3,
			Status:     model.TaskStatusPending,
			Type:       model.TaskTypeSharedPool,
			UserID:     "test_user_id",
			SwapAmount: 10.0,
			CreatedAt:  createTime,
		},
	}

	rewardMap := map[int]*model.RewardRecord{
		1: {
			ID:     1,
			UserID: "test_user_id",
			TaskID: 1,
			Points: 10.0,
		},

		2: {
			ID:     2,
			UserID: "test_user_id",
			TaskID: 2,
			Points: 100.0,
		},
		3: {
			ID:     3,
			UserID: "test_user_id",
			TaskID: 3,
			Points: 55.0,
		},
	}

	t.Run("SearchTasks", func(t *testing.T) {
		testSuite.setUp(t)

		testResponseWriter := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(testResponseWriter)
		startTimeStr := "2021-01-01T00:00:00Z"
		endTimeStr := "2021-01-02T00:00:00Z"
		testContext.Request = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/tasks?user_address=%s&start_time=%s&end_time=%s",
				"test_user_id",
				startTimeStr,
				endTimeStr,
			), nil)

		startTime, _ := time.Parse(time.RFC3339, startTimeStr)
		endTime, _ := time.Parse(time.RFC3339, endTimeStr)

		testSuite.mockedTaskService.EXPECT().
			SearchTasks(&repository.SearchTasksCondition{
				UserID:    "test_user_id",
				StartTime: startTime,
				EndTime:   endTime,
			}).
			Return(&tasks, nil)

		for taskID, record := range rewardMap {
			testSuite.mockedRewardService.EXPECT().
				GetRewardHistoryByTaskID(taskID).
				Return(record, nil)
		}

		testSuite.taskController.SearchTasks(testContext)

		assert.Equal(t, http.StatusOK, testContext.Writer.Status())

		var tasksFromRes response.TaskCollection
		err := json.Unmarshal(testResponseWriter.Body.Bytes(), &tasksFromRes)

		expectedRes := response.TaskCollection{
			{
				Status:            string(model.TaskStatusPending),
				Type:              string(model.TaskTypeSharedPool),
				User:              "test_user_id",
				SwapAmount:        10.0,
				DistributedPoints: 10.0,
				CreatedAt:         tasks[0].CreatedAt,
			},
			{
				Status:            string(model.TaskStatusDone),
				Type:              string(model.TaskTypeOnboarding),
				User:              "test_user_id",
				SwapAmount:        10000.0,
				DistributedPoints: 100.0,
				CreatedAt:         tasks[1].CreatedAt,
			},
			{
				Status:            string(model.TaskStatusPending),
				Type:              string(model.TaskTypeSharedPool),
				User:              "test_user_id",
				SwapAmount:        10.0,
				DistributedPoints: 55.0,
				CreatedAt:         tasks[2].CreatedAt,
			},
		}
		assert.Nil(t, err)
		assert.Equal(t, expectedRes, tasksFromRes)
	})

	t.Run("Get Reward with error", func(t *testing.T) {
		testSuite.setUp(t)

		testResponseWriter := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(testResponseWriter)
		startTimeStr := "2021-01-01T00:00:00Z"
		endTimeStr := "2021-01-02T00:00:00Z"
		testContext.Request = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/tasks?user_address=%s&start_time=%s&end_time=%s",
				"test_user_id",
				startTimeStr,
				endTimeStr,
			), nil)

		startTime, _ := time.Parse(time.RFC3339, startTimeStr)
		endTime, _ := time.Parse(time.RFC3339, endTimeStr)

		testSuite.mockedTaskService.EXPECT().
			SearchTasks(&repository.SearchTasksCondition{
				UserID:    "test_user_id",
				StartTime: startTime,
				EndTime:   endTime,
			}).
			Return(&tasks, nil)

		for taskID, record := range rewardMap {
			if taskID == 1 {
				testSuite.mockedRewardService.EXPECT().
					GetRewardHistoryByTaskID(taskID).
					Return(nil, assert.AnError)
			} else {
				testSuite.mockedRewardService.EXPECT().
					GetRewardHistoryByTaskID(taskID).
					Return(record, nil)
			}
		}

		testSuite.taskController.SearchTasks(testContext)

		assert.Equal(t, http.StatusOK, testContext.Writer.Status())

		var tasksFromRes response.TaskCollection
		err := json.Unmarshal(testResponseWriter.Body.Bytes(), &tasksFromRes)

		expectedRes := response.TaskCollection{
			{
				Status:            string(model.TaskStatusPending),
				Type:              string(model.TaskTypeSharedPool),
				User:              "test_user_id",
				SwapAmount:        10.0,
				DistributedPoints: 0.0,
				CreatedAt:         tasks[0].CreatedAt,
			},
			{
				Status:            string(model.TaskStatusDone),
				Type:              string(model.TaskTypeOnboarding),
				User:              "test_user_id",
				SwapAmount:        10000.0,
				DistributedPoints: 100.0,
				CreatedAt:         tasks[1].CreatedAt,
			},
			{
				Status:            string(model.TaskStatusPending),
				Type:              string(model.TaskTypeSharedPool),
				User:              "test_user_id",
				SwapAmount:        10.0,
				DistributedPoints: 55.0,
				CreatedAt:         tasks[2].CreatedAt,
			},
		}
		assert.Nil(t, err)
		assert.Equal(t, expectedRes, tasksFromRes)
	})

	t.Run("SearchTasks with error", func(t *testing.T) {
		testSuite.setUp(t)

		testResponseWriter := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(testResponseWriter)
		startTimeStr := "2021-01-01T00:00:00Z"
		endTimeStr := "2021-01-02T00:00:00Z"
		testContext.Request = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/tasks?user_address=%s&start_time=%s&end_time=%s",
				"test_user_id",
				startTimeStr,
				endTimeStr,
			), nil)

		startTime, _ := time.Parse(time.RFC3339, startTimeStr)
		endTime, _ := time.Parse(time.RFC3339, endTimeStr)

		testSuite.mockedTaskService.EXPECT().
			SearchTasks(&repository.SearchTasksCondition{
				UserID:    "test_user_id",
				StartTime: startTime,
				EndTime:   endTime,
			}).
			Return(nil, assert.AnError)

		testSuite.taskController.SearchTasks(testContext)

		assert.Equal(t, http.StatusInternalServerError, testContext.Writer.Status())

		var exception map[string]string
		err := json.Unmarshal(testResponseWriter.Body.Bytes(), &exception)

		assert.Nil(t, err)
		assert.Equal(t, assert.AnError.Error(), exception["exception"])
	})

	t.Run("Return empty array when no task found", func(t *testing.T) {
		testSuite.setUp(t)

		testResponseWriter := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(testResponseWriter)
		startTimeStr := "2021-01-01T00:00:00Z"
		endTimeStr := "2021-01-02T00:00:00Z"
		testContext.Request = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/tasks?user_address=%s&start_time=%s&end_time=%s",
				"test_user_id",
				startTimeStr,
				endTimeStr,
			), nil)

		startTime, _ := time.Parse(time.RFC3339, startTimeStr)
		endTime, _ := time.Parse(time.RFC3339, endTimeStr)

		testSuite.mockedTaskService.EXPECT().
			SearchTasks(&repository.SearchTasksCondition{
				UserID:    "test_user_id",
				StartTime: startTime,
				EndTime:   endTime,
			}).
			Return(&[]*model.Task{}, nil)

		testSuite.taskController.SearchTasks(testContext)

		assert.Equal(t, http.StatusOK, testContext.Writer.Status())

		var tasksFromRes response.TaskCollection
		err := json.Unmarshal(testResponseWriter.Body.Bytes(), &tasksFromRes)

		expectedRes := response.TaskCollection{}
		assert.Nil(t, err)
		assert.Equal(t, expectedRes, tasksFromRes)
	})
}
