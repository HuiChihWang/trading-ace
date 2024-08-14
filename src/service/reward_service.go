package service

import (
	"errors"
	"time"
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
		userService:            NewUserService(),
	}
}

func (r *rewardServiceImpl) RewardUser(userID string, TaskID int, points float64) error {
	if points <= 0 {
		return errors.New("points should be greater than 0")
	}

	user, err := r.userService.GetUserByID(userID)

	if err != nil {
		return err
	}

	originalPoints := user.Points

	err = r.userService.UpdateUserPoints(userID, points)
	if err != nil {
		return err
	}

	user, err = r.userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	updatedPoints := user.Points

	rewardRecord := &model.RewardRecord{
		UserID:        userID,
		Points:        points,
		TaskID:        TaskID,
		OriginPoints:  originalPoints,
		UpdatedPoints: updatedPoints,
		CreatedAt:     time.Now().UTC(),
	}
	_, err = r.rewardRecordRepository.CreateRewardRecord(rewardRecord)

	return err
}

func (r *rewardServiceImpl) GetRewardHistory(userID string) ([]*model.RewardRecord, error) {
	return r.rewardRecordRepository.SearchRewardRecords(&repository.RewardRecordSearchCondition{
		UserID: userID,
	})
}
