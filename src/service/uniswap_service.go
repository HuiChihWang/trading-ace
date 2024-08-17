package service

import (
	"errors"
	"fmt"
	"log"
	"time"
	"trading-ace/src/exception"
	"trading-ace/src/model"
	"trading-ace/src/repository"
)

const (
	onboardingAmount      = 1000.0
	onboardingReward      = 100.0
	sharedPoolTotalReward = 10000.0
)

type UniSwapService interface {
	ProcessUniSwapTransaction(senderID string, swapAmount float64) error
	ProcessSharedPool(from time.Time, to time.Time) error
}

type uniSwapServiceImpl struct {
	taskService   TaskService
	userService   UserService
	rewardService RewardService
}

func NewUniSwapService() UniSwapService {
	return &uniSwapServiceImpl{
		taskService:   NewTaskService(),
		userService:   NewUserService(),
		rewardService: NewRewardService(),
	}
}

func (s *uniSwapServiceImpl) ProcessUniSwapTransaction(senderID string, swapAmount float64) error {
	sender, err := s.userService.GetUserByID(senderID)

	if err != nil && !errors.Is(err, exception.UserNotFoundError) {
		return err
	}

	if sender == nil {
		sender, err = s.userService.CreateUser(senderID)
	}

	if err != nil {
		return err
	}

	if !s.isUserAlreadyOnboard(senderID) {
		err := s.processOnBoarding(senderID, swapAmount)

		if err != nil {
			return err
		}
	}

	_, err = s.taskService.CreateTask(senderID, model.TaskTypeSharedPool, swapAmount)

	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("User %s add %f USD to shared pool", senderID, swapAmount))

	return nil
}

func (s *uniSwapServiceImpl) ProcessSharedPool(from time.Time, to time.Time) error {
	tasks, err := s.taskService.SearchTasks(&repository.SearchTasksCondition{
		StartTime: from,
		EndTime:   to,
		Type:      model.TaskTypeSharedPool,
		Status:    model.TaskStatusPending,
	})

	if err != nil {
		return err
	}

	totalSwapAmount := 0.0
	var filteredTasks []*model.Task
	for _, task := range *tasks {
		if !s.isUserAlreadyOnboard(task.UserID) {
			continue
		}

		totalSwapAmount += task.SwapAmount
		filteredTasks = append(filteredTasks, task)
	}

	for _, task := range filteredTasks {
		rewardAmount := sharedPoolTotalReward * task.SwapAmount / totalSwapAmount
		_ = s.rewardService.RewardUser(task.UserID, task.ID, rewardAmount)
		_ = s.taskService.CompleteTask(task.ID)
	}

	return nil
}

func (s *uniSwapServiceImpl) processOnBoarding(userID string, swapAmount float64) error {
	if swapAmount < onboardingAmount {
		log.Println(fmt.Sprintf("User %s does not meet the onboarding requirement", userID))
		return nil
	}

	log.Println(fmt.Sprintf("User %s satisfy onboarding condition with amount %f", userID, swapAmount))

	task, err := s.taskService.CreateTask(userID, model.TaskTypeOnboarding, swapAmount)

	if err != nil {
		return err
	}

	err = s.rewardService.RewardUser(userID, task.ID, onboardingReward)

	if err != nil {
		return err
	}

	return s.taskService.CompleteTask(task.ID)
}

func (s *uniSwapServiceImpl) isUserAlreadyOnboard(userID string) bool {
	tasks, err := s.taskService.SearchTasks(&repository.SearchTasksCondition{
		UserID: userID,
		Type:   model.TaskTypeOnboarding,
	})
	return err == nil && len(*tasks) > 0
}
