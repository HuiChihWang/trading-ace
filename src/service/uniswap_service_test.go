package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"trading-ace/mock/service"
	"trading-ace/src/exception"
	"trading-ace/src/model"
)

type uniSwapServiceTestSuite struct {
	uniSwapService      *uniSwapServiceImpl
	mockedUserService   *service.MockUserService
	mockedTaskService   *service.MockTaskService
	mockedRewardService *service.MockRewardService
}

func (s *uniSwapServiceTestSuite) setUp(t *testing.T) {
	s.mockedUserService = service.NewMockUserService(t)
	s.mockedTaskService = service.NewMockTaskService(t)
	s.mockedRewardService = service.NewMockRewardService(t)
	s.uniSwapService = &uniSwapServiceImpl{
		userService:   s.mockedUserService,
		taskService:   s.mockedTaskService,
		rewardService: s.mockedRewardService,
	}
}

func (s *uniSwapServiceTestSuite) createSharedPoolTasks(taskSetting struct {
	userID     string
	status     model.TaskStatus
	swapAmount float64
	createdAt  time.Time
}) *model.Task {
	return &model.Task{
		UserID:     taskSetting.userID,
		Status:     taskSetting.status,
		SwapAmount: taskSetting.swapAmount,
		Type:       model.TaskTypeSharedPool,
		CreatedAt:  taskSetting.createdAt,
	}
}

var uniSwapTestSuite = &uniSwapServiceTestSuite{}

func TestIsUserAlreadyOnboard(t *testing.T) {
	t.Run("User Already Onboarded", func(t *testing.T) {
		uniSwapTestSuite.setUp(t)

		uniSwapTestSuite.mockedTaskService.EXPECT().GetTasksByUserIDAndType("test_user_address", model.TaskTypeOnboarding).Return([]*model.Task{
			{
				ID:         1,
				UserID:     "test_user_address",
				Type:       model.TaskTypeOnboarding,
				Status:     model.TaskStatusDone,
				SwapAmount: 10.0,
			},
		}, nil).Times(1)

		result := uniSwapTestSuite.uniSwapService.isUserAlreadyOnboard("test_user_address")
		assert.True(t, result)
	})

	t.Run("User Not Onboarded", func(t *testing.T) {
		uniSwapTestSuite.setUp(t)

		uniSwapTestSuite.mockedTaskService.EXPECT().GetTasksByUserIDAndType("test_user_address", model.TaskTypeOnboarding).Return([]*model.Task{}, nil).Times(1)

		result := uniSwapTestSuite.uniSwapService.isUserAlreadyOnboard("test_user_address")
		assert.False(t, result)
	})

	t.Run("Query Error", func(t *testing.T) {
		uniSwapTestSuite.setUp(t)

		uniSwapTestSuite.mockedTaskService.EXPECT().GetTasksByUserIDAndType("test_user_address", model.TaskTypeOnboarding).Return(nil, assert.AnError).Times(1)

		result := uniSwapTestSuite.uniSwapService.isUserAlreadyOnboard("test_user_address")
		assert.False(t, result)
	})
}

func TestUniSwapServiceImpl_ProcessUniSwapTransaction(t *testing.T) {
	t.Run("First Onboard Success and create user", func(t *testing.T) {
		uniSwapTestSuite.setUp(t)

		uniSwapTestSuite.mockedUserService.EXPECT().GetUserByID("test_user_address").Return(nil, exception.UserNotFoundError).Times(1)
		uniSwapTestSuite.mockedUserService.EXPECT().CreateUser("test_user_address").Return(&model.User{
			ID:     "test_user_address",
			Points: 0,
		}, nil).Times(1)

		uniSwapTestSuite.mockedTaskService.EXPECT().GetTasksByUserIDAndType(
			"test_user_address",
			model.TaskTypeOnboarding,
		).Return([]*model.Task{}, nil).Times(1)

		uniSwapTestSuite.mockedTaskService.EXPECT().CreateTask(
			"test_user_address",
			model.TaskTypeOnboarding,
			10000.0,
		).Return(&model.Task{
			ID:         10,
			Type:       model.TaskTypeOnboarding,
			Status:     model.TaskStatusPending,
			SwapAmount: 10000.0,
		}, nil).Times(1)

		uniSwapTestSuite.mockedRewardService.EXPECT().RewardUser("test_user_address", 10, 100.0).Return(nil).Times(1)
		uniSwapTestSuite.mockedTaskService.EXPECT().CompleteTask(10).Return(nil).Times(1)
		uniSwapTestSuite.mockedTaskService.EXPECT().CreateTask("test_user_address", model.TaskTypeSharedPool, 10000.0).Return(&model.Task{}, nil).Times(1)

		err := uniSwapTestSuite.uniSwapService.ProcessUniSwapTransaction("test_user_address", 10000.0)
		assert.Nil(t, err)
	})

	t.Run("First Onboard But have no sufficient amount", func(t *testing.T) {
		uniSwapTestSuite.setUp(t)

		uniSwapTestSuite.mockedUserService.EXPECT().GetUserByID("test_user_address").Return(&model.User{
			ID:     "test_user_address",
			Points: 0,
		}, nil).Times(1)

		uniSwapTestSuite.mockedTaskService.EXPECT().GetTasksByUserIDAndType("test_user_address", model.TaskTypeOnboarding).Return([]*model.Task{}, nil).Times(1)

		err := uniSwapTestSuite.uniSwapService.ProcessUniSwapTransaction("test_user_address", 50.0)
		assert.NotNil(t, err)
	})

	t.Run("User already onboarded", func(t *testing.T) {
		uniSwapTestSuite.setUp(t)

		uniSwapTestSuite.mockedUserService.EXPECT().GetUserByID("test_user_address").Return(&model.User{
			ID:     "test_user_address",
			Points: 0,
		}, nil).Times(1)

		uniSwapTestSuite.mockedTaskService.EXPECT().GetTasksByUserIDAndType(
			"test_user_address",
			model.TaskTypeOnboarding,
		).Return([]*model.Task{
			{
				ID:         1,
				UserID:     "test_user_address",
				Type:       model.TaskTypeOnboarding,
				Status:     model.TaskStatusDone,
				SwapAmount: 10000.0,
			},
		}, nil).Times(1)

		uniSwapTestSuite.mockedTaskService.EXPECT().CreateTask(
			"test_user_address",
			model.TaskTypeSharedPool,
			10000.0,
		).Return(&model.Task{
			ID:         10,
			Type:       model.TaskTypeSharedPool,
			Status:     model.TaskStatusPending,
			SwapAmount: 10000.0,
		}, nil).Times(1)

		err := uniSwapTestSuite.uniSwapService.ProcessUniSwapTransaction("test_user_address", 10000.0)
		assert.Nil(t, err)
	})
}

func TestUniSwapServiceImpl_ProcessSharedPool(t *testing.T) {
	parseTime := func(timeStr string) time.Time {
		t, _ := time.Parse("2006-01-02", timeStr)
		return t
	}

	createdTime := parseTime("2021-01-01").Add(time.Hour * 5)
	tasksPool := []*model.Task{
		{
			ID:         1,
			UserID:     "test_user_1",
			Type:       model.TaskTypeSharedPool,
			Status:     model.TaskStatusPending,
			SwapAmount: 10.0,
			CreatedAt:  createdTime,
		},
		{
			ID:         2,
			UserID:     "test_user_1",
			Type:       model.TaskTypeSharedPool,
			Status:     model.TaskStatusPending,
			SwapAmount: 10.0,
			CreatedAt:  createdTime,
		},
		{
			ID:         3,
			UserID:     "test_user_1",
			Type:       model.TaskTypeSharedPool,
			Status:     model.TaskStatusPending,
			SwapAmount: 20.0,
			CreatedAt:  createdTime,
		},
		{
			ID:         4,
			UserID:     "test_user_2",
			Type:       model.TaskTypeSharedPool,
			Status:     model.TaskStatusPending,
			SwapAmount: 10.0,
			CreatedAt:  createdTime,
		},
		{
			ID:         5,
			UserID:     "test_user_2",
			Type:       model.TaskTypeOnboarding,
			Status:     model.TaskStatusDone,
			SwapAmount: 10.0,
			CreatedAt:  createdTime,
		},
		{
			ID:         5,
			UserID:     "test_user_3",
			Type:       model.TaskTypeSharedPool,
			Status:     model.TaskStatusPending,
			SwapAmount: 10.0,
			CreatedAt:  createdTime,
		},
		{
			ID:         6,
			UserID:     "test_user_2",
			Type:       model.TaskTypeSharedPool,
			Status:     model.TaskStatusDone,
			SwapAmount: 10.0,
			CreatedAt:  createdTime,
		},
	}

	t.Run("Success", func(t *testing.T) {
		uniSwapTestSuite.setUp(t)

		fromTime := parseTime("2021-01-01")
		toTime := parseTime("2021-01-02")

		uniSwapTestSuite.mockedTaskService.EXPECT().GetTasksByDateRange(fromTime, toTime).Return(tasksPool, nil).Times(1)

		IsUsersOnBoarded := map[string]bool{
			"test_user_1": true,
			"test_user_2": true,
			"test_user_3": false,
		}
		for userID, isOnboarded := range IsUsersOnBoarded {
			if isOnboarded {
				uniSwapTestSuite.mockedTaskService.EXPECT().GetTasksByUserIDAndType(userID, model.TaskTypeOnboarding).Return([]*model.Task{
					{
						UserID: userID,
						Type:   model.TaskTypeOnboarding,
						Status: model.TaskStatusDone,
					},
				}, nil)
			} else {
				uniSwapTestSuite.mockedTaskService.EXPECT().GetTasksByUserIDAndType(userID, model.TaskTypeOnboarding).Return([]*model.Task{}, nil)
			}
		}

		filteredTasksAndExpectedPoints := map[*model.Task]float64{
			tasksPool[0]: 2000,
			tasksPool[1]: 2000,
			tasksPool[2]: 4000,
			tasksPool[3]: 2000,
		}

		for task, expectedPoints := range filteredTasksAndExpectedPoints {
			uniSwapTestSuite.mockedRewardService.EXPECT().RewardUser(task.UserID, task.ID, expectedPoints).Return(nil).Times(1)
			uniSwapTestSuite.mockedTaskService.EXPECT().CompleteTask(task.ID).Return(nil).Times(1)
		}

		err := uniSwapTestSuite.uniSwapService.ProcessSharedPool(fromTime, toTime)
		assert.Nil(t, err)
	})

	t.Run("Query Error", func(t *testing.T) {
		uniSwapTestSuite.setUp(t)

		fromTime := parseTime("2021-01-01")
		toTime := parseTime("2021-01-02")
		uniSwapTestSuite.mockedTaskService.EXPECT().GetTasksByDateRange(fromTime, toTime).Return(nil, assert.AnError).Times(1)

		err := uniSwapTestSuite.uniSwapService.ProcessSharedPool(fromTime, toTime)
		assert.NotNil(t, err)
	})
}
