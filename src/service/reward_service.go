package service

import (
	"errors"
	"time"
	"trading-ace/src/model"
	"trading-ace/src/repository"
)

const (
	maxQueryRewardHistoryDuration = 30 * 24 * time.Hour
)

type RewardService interface {
	RewardUser(userID string, TaskID int, points float64) error
	GetRewardHistory(userID string, startTime time.Time, duration time.Duration) ([]*model.RewardRecord, error)
	GetRewardHistoryByTaskID(taskID int) (*model.RewardRecord, error)
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

func (r *rewardServiceImpl) GetRewardHistory(userID string, startTime time.Time, duration time.Duration) ([]*model.RewardRecord, error) {
	if duration <= 0 {
		return nil, errors.New("duration should be greater than 0")
	}

	if duration > maxQueryRewardHistoryDuration {
		return nil, errors.New("duration should be less than 30 days")
	}

	if startTime.After(time.Now().UTC()) {
		return []*model.RewardRecord{}, nil
	}

	return r.rewardRecordRepository.SearchRewardRecords(&repository.RewardRecordSearchCondition{
		UserID:    userID,
		StartTime: startTime,
		Duration:  duration,
	})
}

func (r *rewardServiceImpl) GetRewardHistoryByTaskID(taskID int) (*model.RewardRecord, error) {
	records, err := r.rewardRecordRepository.SearchRewardRecords(&repository.RewardRecordSearchCondition{
		TaskID: taskID,
	})

	if err != nil {
		return nil, err
	}

	return records[0], nil
}
