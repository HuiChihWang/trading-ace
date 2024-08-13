package service

import (
	"database/sql"
	"errors"
	"time"
	"trading-ace/src/model"
	"trading-ace/src/repository"
)

const (
	onboardingAmount      = 1000.0
	onboardingReward      = 100.0
	sharedPoolTotalReward = 10000.0
)

type TaskService interface {
	GetTasksOfUser(userId string) ([]*model.Task, error)
	GetTasksByDateRange(from time.Time, to time.Time) ([]*model.Task, error)
	CreateTask(userId string, taskType model.TaskType, swapAmount float64) (*model.Task, error)
	ProcessOnBoarding(userID string, swapAmount float64) error
	ProcessSharedPool(from time.Time, to time.Time) error
	IsUserOnboardingCompleted(userID string) bool
}

type taskServiceImpl struct {
	userService    UserService
	rewardService  RewardService
	taskRepository repository.TaskRepository
}

func NewTaskService() TaskService {
	return &taskServiceImpl{
		userService:    NewUserService(),
		rewardService:  NewRewardService(),
		taskRepository: repository.NewTaskRepository(),
	}
}

func (s *taskServiceImpl) GetTasksOfUser(userId string) ([]*model.Task, error) {
	return s.taskRepository.GetTasksByUserID(userId)
}

func (s *taskServiceImpl) GetTasksByDateRange(from time.Time, to time.Time) ([]*model.Task, error) {
	return s.taskRepository.GetTasksByDateRange(from, to)
}

func (s *taskServiceImpl) CreateTask(userId string, taskType model.TaskType, swapAmount float64) (*model.Task, error) {
	task := model.NewTask(userId, taskType, swapAmount)
	return s.taskRepository.CreateTask(task)
}

func (s *taskServiceImpl) IsUserOnboardingCompleted(userID string) bool {
	onboardingTasks, err := s.taskRepository.GetTasksByUserIDAndType(userID, model.TaskTypeOnboarding)

	if err != nil {
		return false
	}

	return len(onboardingTasks) > 0
}

func (s *taskServiceImpl) ProcessOnBoarding(userID string, swapAmount float64) error {
	if swapAmount < onboardingAmount {
		return errors.New("swap amount does not meet the requirement")
	}

	if s.IsUserOnboardingCompleted(userID) {
		return errors.New("onboarding already completed")
	}

	task, err := s.CreateTask(userID, model.TaskTypeOnboarding, swapAmount)

	if err != nil {
		return err
	}

	err = s.rewardService.RewardUser(userID, task.ID, onboardingReward)

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

func (s *taskServiceImpl) ProcessSharedPool(from time.Time, to time.Time) error {
	tasks, err := s.taskRepository.GetTasksByDateRange(from, to)

	if err != nil {
		return err
	}

	totalSwapAmount := 0.0
	filteredTasks := make([]*model.Task, 0)
	for _, task := range tasks {
		if task.Type != model.TaskTypeSharedPool || task.Status == model.TaskStatusDone {
			totalSwapAmount += task.SwapAmount
			filteredTasks = append(filteredTasks, task)
		}
	}

	for _, task := range filteredTasks {
		rewardAmount := sharedPoolTotalReward * task.SwapAmount / totalSwapAmount
		err := s.rewardService.RewardUser(task.UserID, task.ID, rewardAmount)
		if err != nil {
			return err
		}

		task.Status = model.TaskStatusDone
		task.CompletedAt = sql.NullTime{
			Time:  time.Now().In(time.UTC),
			Valid: true,
		}

		_, err = s.taskRepository.UpdateTask(task)

		if err != nil {
			return err
		}
	}

	for _, task := range tasks {
		task.Status = model.TaskStatusDone
		task.CompletedAt = sql.NullTime{
			Time:  time.Now().In(time.UTC),
			Valid: true,
		}
		_, err := s.taskRepository.UpdateTask(task)
		if err != nil {
			return err
		}
	}

	return nil
}
