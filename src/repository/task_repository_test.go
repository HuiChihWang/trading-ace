package repository

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"trading-ace/src/database"
	"trading-ace/src/model"
)

func TestTaskRepositoryImpl(t *testing.T) {
	setUpTaskRepo := func(t *testing.T) *taskRepositoryImpl {
		dbInstance := database.GetDBInstance()

		t.Cleanup(func() {
			dbInstance.Exec("DELETE FROM tasks")
		})

		return &taskRepositoryImpl{
			dbInstance: dbInstance,
		}
	}

	t.Run("CreateTask", func(t *testing.T) {
		taskRepo := setUpTaskRepo(t)
		task := model.NewTask("test_user_id", model.TaskTypeOnboarding, 50)
		createdTask, err := taskRepo.CreateTask(task)

		if err != nil {
			t.Errorf("CreateTask() error = %v", err)
		}

		assert.NotEmpty(t, task.ID)
		assert.Equal(t, "test_user_id", createdTask.UserID)
		assert.Equal(t, model.TaskTypeOnboarding, createdTask.Type)
		assert.Equal(t, model.TaskStatusPending, createdTask.Status)
		assert.Equal(t, task.CreatedAt, createdTask.CreatedAt)
		assert.Equal(t, 50.0, createdTask.SwapAmount)
		assert.Equal(t, createdTask.CompletedAt, sql.NullTime{})
	})

	t.Run("UpdateTask", func(t *testing.T) {
		taskRepo := setUpTaskRepo(t)
		task := model.NewTask("test_user_id", model.TaskTypeOnboarding, 50)
		task, _ = taskRepo.CreateTask(task)

		task.Status = model.TaskStatusDone
		task.CompletedAt = sql.NullTime{
			Time:  time.Now().In(time.UTC),
			Valid: true,
		}
		updatedTask, err := taskRepo.UpdateTask(task)

		if err != nil {
			t.Errorf("UpdateTask() error = %v", err)
		}

		assert.Equal(t, task.ID, updatedTask.ID)
		assert.Equal(t, task.UserID, updatedTask.UserID)
		assert.Equal(t, task.Type, updatedTask.Type)
		assert.Equal(t, task.SwapAmount, updatedTask.SwapAmount)
		assert.Equal(t, model.TaskStatusDone, updatedTask.Status)
		assert.NotEmpty(t, updatedTask.CreatedAt)
		assert.True(t, time.Now().Sub(updatedTask.CompletedAt.Time) < time.Second)
	})

	t.Run("SearchTasks", func(t *testing.T) {
		taskRepo := setUpTaskRepo(t)

		tasksToBeInsert := []*model.Task{
			{
				UserID:     "user_1",
				Type:       model.TaskTypeOnboarding,
				Status:     model.TaskStatusPending,
				SwapAmount: 50,
				CreatedAt:  time.Now().Add(-time.Hour * 10),
			},
			{
				UserID:     "user_1",
				Type:       model.TaskTypeSharedPool,
				Status:     model.TaskStatusDone,
				SwapAmount: 100,
				CreatedAt:  time.Now().Add(-time.Hour * 9),
			},
			{
				UserID:     "user_2",
				Type:       model.TaskTypeOnboarding,
				Status:     model.TaskStatusPending,
				SwapAmount: 75,
				CreatedAt:  time.Now().Add(-time.Hour * 8),
			},
			{
				UserID:     "user_2",
				Type:       model.TaskTypeSharedPool,
				Status:     model.TaskStatusDone,
				SwapAmount: 150,
				CreatedAt:  time.Now().Add(-time.Hour * 7),
			},
			{
				UserID:     "user_3",
				Type:       model.TaskTypeOnboarding,
				Status:     model.TaskStatusPending,
				SwapAmount: 200,
				CreatedAt:  time.Now().Add(-time.Hour * 6),
			},
			{
				UserID:     "user_3",
				Type:       model.TaskTypeSharedPool,
				Status:     model.TaskStatusDone,
				SwapAmount: 250,
				CreatedAt:  time.Now().Add(-time.Hour * 5),
			},
			{
				UserID:     "user_1",
				Type:       model.TaskTypeSharedPool,
				Status:     model.TaskStatusPending,
				SwapAmount: 300,
				CreatedAt:  time.Now().Add(-time.Hour * 4),
			},
			{
				UserID:     "user_2",
				Type:       model.TaskTypeSharedPool,
				Status:     model.TaskStatusDone,
				SwapAmount: 350,
				CreatedAt:  time.Now().Add(-time.Hour * 3),
			},
			{
				UserID:     "user_3",
				Type:       model.TaskTypeSharedPool,
				Status:     model.TaskStatusPending,
				SwapAmount: 400,
				CreatedAt:  time.Now().Add(-time.Hour * 2),
			},
			{
				UserID:     "user_4",
				Type:       model.TaskTypeOnboarding,
				Status:     model.TaskStatusDone,
				SwapAmount: 450,
				CreatedAt:  time.Now().Add(-time.Hour * 1),
			},
		}

		for _, task := range tasksToBeInsert {
			newTask, _ := taskRepo.CreateTask(task)
			task.ID = newTask.ID
		}

		t.Run("Search By User", func(t *testing.T) {
			tasks, err := taskRepo.SearchTasks(&SearchTasksCondition{
				UserID: "user_1",
			})
			assert.NoError(t, err)
			assert.Equal(t, 3, len(tasks))

			expectedTaskIndexes := []int{6, 1, 0}
			for i, task := range tasks {
				assert.Equal(t, tasksToBeInsert[expectedTaskIndexes[i]].ID, task.ID)
			}
		})

		t.Run("Search By Type", func(t *testing.T) {
			tasks, err := taskRepo.SearchTasks(&SearchTasksCondition{
				Type: model.TaskTypeOnboarding,
			})
			assert.NoError(t, err)
			assert.Equal(t, 4, len(tasks))

			expectedTaskIndexes := []int{9, 4, 2, 0}
			for i, task := range tasks {
				assert.Equal(t, tasksToBeInsert[expectedTaskIndexes[i]].ID, task.ID)
			}
		})

		t.Run("Search By Status", func(t *testing.T) {
			tasks, err := taskRepo.SearchTasks(&SearchTasksCondition{
				Status: model.TaskStatusDone,
			})
			assert.NoError(t, err)
			assert.Equal(t, 5, len(tasks))

			expectedTaskIndexes := []int{9, 7, 5, 3, 1}
			for i, task := range tasks {
				assert.Equal(t, tasksToBeInsert[expectedTaskIndexes[i]].ID, task.ID)
			}
		})

		t.Run("Search, Type, Status", func(t *testing.T) {
			tasks, err := taskRepo.SearchTasks(&SearchTasksCondition{
				Type:   model.TaskTypeOnboarding,
				Status: model.TaskStatusDone,
			})
			assert.NoError(t, err)
			assert.Equal(t, 1, len(tasks))

			expectedTaskIndexes := []int{9}
			for i, task := range tasks {
				assert.Equal(t, tasksToBeInsert[expectedTaskIndexes[i]].ID, task.ID)
			}
		})

		t.Run("Search By Time Range", func(t *testing.T) {
			tasks, err := taskRepo.SearchTasks(&SearchTasksCondition{
				StartTime: time.Now().Add(-time.Hour*3 - time.Minute),
				EndTime:   time.Now(),
			})
			assert.NoError(t, err)
			assert.Equal(t, 3, len(tasks))

			expectedTaskIndexes := []int{9, 8, 7}
			for i, task := range tasks {
				assert.Equal(t, tasksToBeInsert[expectedTaskIndexes[i]].ID, task.ID)
			}
		})

		t.Run("InValid Search Range", func(t *testing.T) {
			testConditions := map[string][]time.Time{
				"Start Time > End Time": {time.Now().Add(time.Hour), time.Now()},
				"Empty Start Time":      {time.Time{}, time.Now()},
				"Empty End Time":        {time.Now(), time.Time{}},
			}

			for _, timeRange := range testConditions {
				tasks, err := taskRepo.SearchTasks(&SearchTasksCondition{
					StartTime: timeRange[0],
					EndTime:   timeRange[1],
				})
				assert.NotNil(t, err)
				assert.Nil(t, tasks)
			}
		})
	})
}
