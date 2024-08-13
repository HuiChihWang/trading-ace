package service

import (
	"database/sql"
	"time"
	"trading-ace/src/model"
	"trading-ace/src/repository"
)

type TaskService interface {
	GetTasksOfUser(userId string) ([]*model.Task, error)
	GetTasksByDateRange(from time.Time, to time.Time) ([]*model.Task, error)
	GetTasksByUserIDAndType(userID string, taskType model.TaskType) ([]*model.Task, error)
	CreateTask(userId string, taskType model.TaskType, swapAmount float64) (*model.Task, error)
	CompleteTask(taskID int) error
}

type taskServiceImpl struct {
	taskRepository repository.TaskRepository
}

func NewTaskService() TaskService {
	return &taskServiceImpl{
		taskRepository: repository.NewTaskRepository(),
	}
}

func (s *taskServiceImpl) GetTasksOfUser(userId string) ([]*model.Task, error) {
	return s.taskRepository.GetTasksByUserID(userId)
}

func (s *taskServiceImpl) GetTasksByDateRange(from time.Time, to time.Time) ([]*model.Task, error) {
	return s.taskRepository.GetTasksByDateRange(from, to)
}

func (s *taskServiceImpl) GetTasksByUserIDAndType(userID string, taskType model.TaskType) ([]*model.Task, error) {
	return s.taskRepository.GetTasksByUserIDAndType(userID, taskType)
}

func (s *taskServiceImpl) CreateTask(userId string, taskType model.TaskType, swapAmount float64) (*model.Task, error) {
	task := model.NewTask(userId, taskType, swapAmount)
	return s.taskRepository.CreateTask(task)
}

func (s *taskServiceImpl) CompleteTask(taskID int) error {
	task, err := s.taskRepository.GetTaskByID(taskID)

	if err != nil {
		return err
	}

	task.Status = model.TaskStatusDone
	task.CompletedAt = sql.NullTime{
		Time:  time.Now().In(time.UTC),
		Valid: true,
	}

	_, err = s.taskRepository.UpdateTask(task)

	return err
}
