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

// test searchTasks
func TestTaskServiceImpl_SearchTasks(t *testing.T) {
	testSuite := &taskServiceTestSuite{}

	tasksFromRepo := []*model.Task{
		{
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
		},
		{
			ID:         2,
			UserID:     "test_user_id",
			Type:       model.TaskTypeSharedPool,
			Status:     model.TaskStatusDone,
			SwapAmount: 100.0,
			CompletedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			CreatedAt: time.Now().Add(-time.Hour * 10),
		},
	}

	searchCond := &realRepo.SearchTasksCondition{
		UserID: "test_user_id",
	}

	t.Run("SearchTasks", func(t *testing.T) {
		testSuite.setUp(t)

		testSuite.mockedTaskRepository.EXPECT().SearchTasks(searchCond).Return(tasksFromRepo, nil).Times(1)

		tasks, err := testSuite.taskService.SearchTasks(searchCond)
		assert.Nil(t, err)
		assert.Equal(t, &tasksFromRepo, tasks)
	})

	t.Run("SearchTasksFail", func(t *testing.T) {
		testSuite.setUp(t)

		testSuite.mockedTaskRepository.EXPECT().SearchTasks(searchCond).Return(
			nil, assert.AnError).Times(1)

		tasks, err := testSuite.taskService.SearchTasks(searchCond)
		assert.Nil(t, tasks)
		assert.NotNil(t, err)
	})
}
