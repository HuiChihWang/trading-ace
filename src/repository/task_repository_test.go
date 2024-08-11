package repository

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"trading-ace/src/config"
	"trading-ace/src/database"
	"trading-ace/src/model"
)

func TestTaskRepositoryImpl(t *testing.T) {
	setUpTaskRepo := func(t *testing.T) *taskRepositoryImpl {
		dbInstance := database.CreateDBInstance(&config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5435",
			Username: "postgres",
			Password: "postgres",
			DBName:   "trading_ace_test",
		})

		t.Cleanup(func() {
			dbInstance.Exec("DELETE FROM tasks")
			fmt.Println("[Tear Down] Cleaned up tasks table")
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

	t.Run("GetTaskByID", func(t *testing.T) {
		taskRepo := setUpTaskRepo(t)
		task := model.NewTask("test_user_id", model.TaskTypeOnboarding, 50)
		task, _ = taskRepo.CreateTask(task)

		fetchedTask, err := taskRepo.GetTaskByID(task.ID)

		if err != nil {
			t.Errorf("GetTaskByID() error = %v", err)
		}

		assert.Equal(t, task.ID, fetchedTask.ID)
		assert.Equal(t, task.UserID, fetchedTask.UserID)
		assert.Equal(t, task.Type, fetchedTask.Type)
		assert.Equal(t, task.SwapAmount, fetchedTask.SwapAmount)
		assert.Equal(t, task.Status, fetchedTask.Status)
		assert.NotEmpty(t, fetchedTask.CreatedAt)
		assert.Equal(t, fetchedTask.CompletedAt, sql.NullTime{})
	})

	t.Run("GetTasksByUserID", func(t *testing.T) {
		taskRepo := setUpTaskRepo(t)
		for i := 0; i < 10; i++ {
			task := model.NewTask("test_user_id", model.TaskTypeOnboarding, 50)
			_, _ = taskRepo.CreateTask(task)
		}

		tasks, err := taskRepo.GetTasksByUserID("test_user_id")

		if err != nil {
			t.Errorf("GetTasksByUserID() error = %v", err)
		}

		assert.Equal(t, 10, len(tasks))
		for _, task := range tasks {
			assert.NotEmpty(t, task.ID)
			assert.Equal(t, "test_user_id", task.UserID)
			assert.Equal(t, model.TaskTypeOnboarding, task.Type)
			assert.Equal(t, model.TaskStatusPending, task.Status)
			assert.Equal(t, 50.0, task.SwapAmount)
			assert.NotEmpty(t, task.CreatedAt)
			assert.Equal(t, task.CompletedAt, sql.NullTime{})
		}
	})

	t.Run("GetTasksByUserIDAndType", func(t *testing.T) {
		taskRepo := setUpTaskRepo(t)
		for i := 0; i < 10; i++ {
			task := model.NewTask("test_user_id", model.TaskTypeOnboarding, 50)
			_, _ = taskRepo.CreateTask(task)
		}

		tasks, err := taskRepo.GetTasksByUserIDAndType("test_user_id", model.TaskTypeOnboarding)

		if err != nil {
			t.Errorf("GetTasksByUserIDAndType() error = %v", err)
		}

		assert.Equal(t, 10, len(tasks))
		for _, task := range tasks {
			assert.Equal(t, "test_user_id", task.UserID)
			assert.Equal(t, model.TaskTypeOnboarding, task.Type)
			assert.Equal(t, model.TaskStatusPending, task.Status)
			assert.Equal(t, 50.0, task.SwapAmount)
			assert.NotEmpty(t, task.CreatedAt)
			assert.NotEmpty(t, task.ID)
			assert.Equal(t, task.CompletedAt, sql.NullTime{})
		}
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
		assert.Equal(t, task.CompletedAt, updatedTask.CompletedAt)
	})

	t.Run("GetTasksByDateRange", func(t *testing.T) {
		taskRepo := setUpTaskRepo(t)

		// create in today
		for i := 0; i < 10; i++ {
			task := model.NewTask("test_user_id", model.TaskTypeOnboarding, 50)
			task.CreatedAt = time.Now().Add(-time.Hour * 12)
			_, _ = taskRepo.CreateTask(task)
		}

		// create in last week
		for i := 0; i < 10; i++ {
			task := model.NewTask("test_user_id", model.TaskTypeOnboarding, 50)
			task.CreatedAt = time.Now().Add(-time.Hour * 24 * 8)
			_, _ = taskRepo.CreateTask(task)
		}

		t.Run("GetTasksByDateRange - Check Today", func(t *testing.T) {
			from := time.Now().Add(-time.Hour * 24)
			to := time.Now()
			tasks, err := taskRepo.GetTasksByDateRange(from, to)

			if err != nil {
				t.Errorf("GetTasksByDateRange() error = %v", err)
			}

			assert.Equal(t, 10, len(tasks))
			for _, task := range tasks {
				assert.NotEmpty(t, task.ID)
				assert.Equal(t, "test_user_id", task.UserID)
				assert.Equal(t, model.TaskTypeOnboarding, task.Type)
				assert.Equal(t, model.TaskStatusPending, task.Status)
				assert.Equal(t, 50.0, task.SwapAmount)
				assert.NotEmpty(t, task.CreatedAt)
				assert.Equal(t, task.CompletedAt, sql.NullTime{})
			}
		})
		t.Run("GetTasksByDateRange - Check Last Week", func(t *testing.T) {
			from := time.Now().Add(-time.Hour * 24 * 14)
			to := time.Now().Add(-time.Hour * 24 * 7)
			tasks, err := taskRepo.GetTasksByDateRange(from, to)

			if err != nil {
				t.Errorf("GetTasksByDateRange() error = %v", err)
			}

			assert.Equal(t, 10, len(tasks))
			for _, task := range tasks {
				assert.NotEmpty(t, task.ID)
				assert.Equal(t, "test_user_id", task.UserID)
				assert.Equal(t, model.TaskTypeOnboarding, task.Type)
				assert.Equal(t, model.TaskStatusPending, task.Status)
				assert.Equal(t, 50.0, task.SwapAmount)
				assert.NotEmpty(t, task.CreatedAt)
				assert.Equal(t, task.CompletedAt, sql.NullTime{})
			}
		})
		t.Run("GetTasksByDateRange - Check Two week before", func(t *testing.T) {
			from := time.Now().Add(-time.Hour * 24 * 21)
			to := time.Now().Add(-time.Hour * 24 * 14)
			tasks, err := taskRepo.GetTasksByDateRange(from, to)

			if err != nil {
				t.Errorf("GetTasksByDateRange() error = %v", err)
			}

			assert.Equal(t, 0, len(tasks))
		})
	})
}
