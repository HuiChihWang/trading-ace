package service

import (
	"trading-ace/src/model"
	"trading-ace/src/repository"
)

type RewardService interface {
	RewardUser(userID string, TaskID int, points float64) error
	GetRewardHistory(userID string) ([]*model.RewardRecord, error)
}

type rewardServiceImpl struct {
	rewardRecordRepository repository.RewardRecordRepository
	userService            UserService
}

func NewRewardService() RewardService {
	return &rewardServiceImpl{
		rewardRecordRepository: repository.NewRewardRecordRepository(),
	}
}

func (r *rewardServiceImpl) RewardUser(userID string, TaskID int, points float64) error {
	err := r.userService.UpdateUserPoints(userID, points)
	if err != nil {
		return err
	}

	_, err = r.rewardRecordRepository.CreateRewardRecord(userID, points, TaskID)
	return err
}

func (r *rewardServiceImpl) GetRewardHistory(userID string) ([]*model.RewardRecord, error) {
	return r.rewardRecordRepository.GetRewardRecordsByUserID(userID)
}
