package service

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
	"trading-ace/mock/repository"
	"trading-ace/src/model"
	realRepo "trading-ace/src/repository"
)

type taskServiceTestSuite struct {
	taskService          TaskService
	mockedTaskRepository *repository.MockTaskRepository
}

func (s *taskServiceTestSuite) setUp(t *testing.T) {
	s.mockedTaskRepository = repository.NewMockTaskRepository(t)
	s.taskService = &taskServiceImpl{
		taskRepository: s.mockedTaskRepository,
	}
}

func TestTaskServiceImpl_GetTasksOfUser(t *testing.T) {
	testSuite := &taskServiceTestSuite{}

	t.Run("GetTasksOfUser", func(t *testing.T) {
		testSuite.setUp(t)

		tasksFromRepo := []*model.Task{
			{
				ID:         1,
				UserID:     "test_user_id",
				Type:       model.TaskTypeOnboarding,
				Status:     model.TaskStatusDone,
				SwapAmount: 10.0,
			},
		}
		testSuite.mockedTaskRepository.EXPECT().SearchTasks(&realRepo.SearchTasksCondition{
			UserID: "test_user_id",
		}).Return(tasksFromRepo, nil).Times(1)

		tasks, err := testSuite.taskService.GetTasksOfUser("test_user_id")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(tasks))
		assert.Equal(t, tasksFromRepo[0], tasks[0])
	})

	t.Run("GetTasksOfUserFail", func(t *testing.T) {
		testSuite.setUp(t)

		testSuite.mockedTaskRepository.EXPECT().SearchTasks(&realRepo.SearchTasksCondition{
			UserID: "test_user_id",
		}).Return(nil, assert.AnError).Times(1)

		tasks, err := testSuite.taskService.GetTasksOfUser("test_user_id")
		assert.Nil(t, tasks)
		assert.NotNil(t, err)
	})
}

func TestTaskServiceImpl_CompleteTask(t *testing.T) {
	testSuite := &taskServiceTestSuite{}

	t.Run("CompleteTask", func(t *testing.T) {
		testSuite.setUp(t)

		taskFromRepo := &model.Task{
			ID:         1,
			UserID:     "test_user_id",
			Type:       model.TaskTypeOnboarding,
			Status:     model.TaskStatusPending,
			SwapAmount: 10.0,
			CompletedAt: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		}
		testSuite.mockedTaskRepository.EXPECT().GetTaskByID(1).Return(taskFromRepo, nil).Times(1)
		testSuite.mockedTaskRepository.EXPECT().UpdateTask(taskFromRepo).Return(nil, nil).Times(1)

		err := testSuite.taskService.CompleteTask(1)
		assert.Nil(t, err)
		assert.Equal(t, model.TaskStatusDone, taskFromRepo.Status)
		assert.True(t, taskFromRepo.CompletedAt.Valid)
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		testSuite.setUp(t)

		testSuite.mockedTaskRepository.EXPECT().GetTaskByID(1).Return(nil, assert.AnError).Times(1)

		err := testSuite.taskService.CompleteTask(1)
		assert.NotNil(t, err)
	})

	t.Run("UpdateTaskFail", func(t *testing.T) {
		testSuite.setUp(t)

		taskFromRepo := &model.Task{
			ID:         1,
			UserID:     "test_user_id",
			Type:       model.TaskTypeOnboarding,
			Status:     model.TaskStatusPending,
			SwapAmount: 10.0,
			CompletedAt: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		}
		testSuite.mockedTaskRepository.EXPECT().GetTaskByID(1).Return(taskFromRepo, nil).Times(1)
		testSuite.mockedTaskRepository.EXPECT().UpdateTask(taskFromRepo).Return(nil, assert.AnError).Times(1)

		err := testSuite.taskService.CompleteTask(1)
		assert.NotNil(t, err)
	})
}

func TestTaskServiceImpl_GetTasksByDateRange(t *testing.T) {
	testSuite := &taskServiceTestSuite{}
	t.Run("GetTasksByDateRange", func(t *testing.T) {
		testSuite.setUp(t)

		from := time.Now().Add(-time.Hour * 12)
		to := time.Now().Add(time.Hour * 12)
		tasksFromRepo := []*model.Task{
			{
				ID:         1,
				UserID:     "test_user_id",
				Type:       model.TaskTypeOnboarding,
				Status:     model.TaskStatusDone,
				SwapAmount: 10.0,
				CreatedAt:  time.Now(),
			},
		}
		testSuite.mockedTaskRepository.EXPECT().SearchTasks(&realRepo.SearchTasksCondition{
			StartTime: from,
			EndTime:   to,
		}).Return(tasksFromRepo, nil).Times(1)

		tasks, err := testSuite.taskService.GetTasksByDateRange(from, to)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(tasks))
		assert.Equal(t, tasksFromRepo[0], tasks[0])
	})

	t.Run("GetTasksByDateRangeFail", func(t *testing.T) {
		testSuite.setUp(t)

		from := time.Now()
		to := time.Now().Add(time.Hour * 24)
		testSuite.mockedTaskRepository.EXPECT().SearchTasks(&realRepo.SearchTasksCondition{
			StartTime: from,
			EndTime:   to,
		}).Return(nil, assert.AnError).Times(1)

		tasks, err := testSuite.taskService.GetTasksByDateRange(from, to)
		assert.Nil(t, tasks)
		assert.NotNil(t, err)
	})
}

func TestTaskServiceImpl_GetTasksByUserIDAndType(t *testing.T) {
	testSuite := &taskServiceTestSuite{}
	t.Run("GetTasksByUserIDAndType", func(t *testing.T) {
		testSuite.setUp(t)

		tasksFromRepo := []*model.Task{
			{
				ID:         1,
				UserID:     "test_user_id",
				Type:       model.TaskTypeOnboarding,
				Status:     model.TaskStatusDone,
				SwapAmount: 10.0,
			},
		}
		testSuite.mockedTaskRepository.EXPECT().SearchTasks(&realRepo.SearchTasksCondition{
			UserID: "test_user_id",
			Type:   model.TaskTypeOnboarding,
		}).Return(tasksFromRepo, nil).Times(1)

		tasks, err := testSuite.taskService.GetTasksByUserIDAndType("test_user_id", model.TaskTypeOnboarding)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(tasks))
		assert.Equal(t, tasksFromRepo[0], tasks[0])
	})

	t.Run("GetTasksByUserIDAndTypeFail", func(t *testing.T) {
		testSuite.setUp(t)

		testSuite.mockedTaskRepository.EXPECT().SearchTasks(&realRepo.SearchTasksCondition{
			UserID: "test_user_id",
			Type:   model.TaskTypeOnboarding,
		}).Return(nil, assert.AnError).Times(1)

		tasks, err := testSuite.taskService.GetTasksByUserIDAndType("test_user_id", model.TaskTypeOnboarding)
		assert.Nil(t, tasks)
		assert.NotNil(t, err)
	})
}

func TestTaskServiceImpl_CreateTask(t *testing.T) {
	testSuite := &taskServiceTestSuite{}

	t.Run("CreateTask", func(t *testing.T) {
		testSuite.setUp(t)

		taskCreatedByRepo := &model.Task{
			ID:         1,
			UserID:     "test_user_id",
			Type:       model.TaskTypeOnboarding,
			Status:     model.TaskStatusPending,
			SwapAmount: 10.0,
			CompletedAt: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
			CreatedAt: time.Now(),
		}

		testSuite.mockedTaskRepository.EXPECT().CreateTask(mock.MatchedBy(
			func(task *model.Task) bool {
				return task.UserID == "test_user_id" &&
					task.Type == model.TaskTypeOnboarding &&
					task.SwapAmount == 10.0 &&
					task.Status == model.TaskStatusPending &&
					task.CompletedAt.Valid == false
			},
		)).Return(taskCreatedByRepo, nil).Times(1)

		newTask, err := testSuite.taskService.CreateTask("test_user_id", model.TaskTypeOnboarding, 10.0)
		assert.Nil(t, err)
		assert.Equal(t, model.TaskTypeOnboarding, newTask.Type)
		assert.Equal(t, 10.0, newTask.SwapAmount)
		assert.Equal(t, "test_user_id", newTask.UserID)
		assert.Equal(t, model.TaskStatusPending, newTask.Status)
		assert.NotEmpty(t, newTask.ID)
		assert.NotEmpty(t, newTask.CreatedAt)
		assert.False(t, newTask.CompletedAt.Valid)
	})

	t.Run("CreateTaskFail", func(t *testing.T) {
		testSuite.setUp(t)

		testSuite.mockedTaskRepository.EXPECT().CreateTask(mock.Anything).Return(
			nil, assert.AnError).Times(1)

		newTask, err := testSuite.taskService.CreateTask("test_user_id", model.TaskTypeOnboarding, 10.0)
		assert.Nil(t, newTask)
		assert.NotNil(t, err)
	})
}
