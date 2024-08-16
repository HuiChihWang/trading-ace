package service

import (
	"database/sql"
	"time"
	"trading-ace/src/model"
	"trading-ace/src/repository"
)

type TaskService interface {
	CreateTask(userId string, taskType model.TaskType, swapAmount float64) (*model.Task, error)
	CompleteTask(taskID int) error
	SearchTasks(condition *repository.SearchTasksCondition) (*[]*model.Task, error)
}

type taskServiceImpl struct {
	taskRepository repository.TaskRepository
}

func NewTaskService() TaskService {
	return &taskServiceImpl{
		taskRepository: repository.NewTaskRepository(),
	}
}

func (s *taskServiceImpl) SearchTasks(condition *repository.SearchTasksCondition) (*[]*model.Task, error) {
	tasks, err := s.taskRepository.SearchTasks(condition)

	if err != nil {
		return nil, err
	}

	return &tasks, nil
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
