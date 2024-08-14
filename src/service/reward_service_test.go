package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"trading-ace/mock/repository"
	"trading-ace/mock/service"
	"trading-ace/src/model"
	repoReal "trading-ace/src/repository"
)

var rewardService RewardService
var mockedUserService *service.MockUserService
var mockedRewardRecordRepository *repository.MockRewardRecordRepository

func setUpRewardService(t *testing.T) {
	mockedUserService = service.NewMockUserService(t)
	mockedRewardRecordRepository = repository.NewMockRewardRecordRepository(t)
	rewardService = &rewardServiceImpl{
		rewardRecordRepository: mockedRewardRecordRepository,
		userService:            mockedUserService,
	}
}

func TestRewardServiceImpl_RewardUser(t *testing.T) {
	t.Run("RewardUser", func(t *testing.T) {
		setUpRewardService(t)

		mockedUserService.EXPECT().GetUserByID("test_user_id").Return(&model.User{
			ID:     "test_user_id",
			Points: 0.0,
		}, nil).Times(1)

		mockedUserService.EXPECT().UpdateUserPoints("test_user_id", 10.0).Return(nil).Times(1)
		mockedRewardRecordRepository.EXPECT().CreateRewardRecord(mock.MatchedBy(
			func(rewardRecord *model.RewardRecord) bool {
				return rewardRecord.UserID == "test_user_id" &&
					rewardRecord.Points == 10.0 &&
					rewardRecord.TaskID == 1 &&
					rewardRecord.OriginPoints == 0.0 &&
					rewardRecord.UpdatedPoints == 10.0
			},
		)).Return(
			&model.RewardRecord{
				UserID:        "test_user_id",
				Points:        10.0,
				TaskID:        1,
				OriginPoints:  0.0,
				UpdatedPoints: 10.0,
			}, nil).Times(1)

		mockedUserService.EXPECT().GetUserByID("test_user_id").Return(&model.User{
			ID:     "test_user_id",
			Points: 10.0,
		}, nil).Times(1)

		err := rewardService.RewardUser("test_user_id", 1, 10.0)
		if err != nil {
			t.Errorf("RewardUser() exception = %v", err)
		}
	})

	t.Run("RewardWithNegativeOrZeroPoints", func(t *testing.T) {
		setUpRewardService(t)

		t.Run("NegativePoints", func(t *testing.T) {
			err := rewardService.RewardUser("test_user_id", 1, -10.0)
			if err == nil {
				t.Errorf("RewardUser() expected error but got nil")
			}
		})

		t.Run("ZeroPoints", func(t *testing.T) {
			err := rewardService.RewardUser("test_user_id", 1, 0.0)
			if err == nil {
				t.Errorf("RewardUser() expected error but got nil")
			}
		})
	})

	t.Run("RewardUserFail", func(t *testing.T) {
		setUpRewardService(t)

		mockedUserService.EXPECT().GetUserByID("test_user_id").Return(&model.User{
			ID:     "test_user_id",
			Points: 0.0,
		}, nil).Times(1)
		mockedUserService.EXPECT().UpdateUserPoints("test_user_id", 10.0).Return(assert.AnError).Times(1)

		err := rewardService.RewardUser("test_user_id", 1, 10.0)
		if err == nil {
			t.Errorf("RewardUser() expected error but got nil")
		}
	})

	t.Run("CreateRewardRecordFail", func(t *testing.T) {
		setUpRewardService(t)

		mockedUserService.EXPECT().GetUserByID("test_user_id").Return(&model.User{
			ID:     "test_user_id",
			Points: 0.0,
		}, nil).Times(1)

		mockedUserService.EXPECT().UpdateUserPoints("test_user_id", 10.0).Return(nil).Times(1)
		mockedUserService.EXPECT().GetUserByID("test_user_id").Return(&model.User{
			ID:     "test_user_id",
			Points: 10.0,
		}, nil).Times(1)

		mockedRewardRecordRepository.EXPECT().CreateRewardRecord(mock.Anything).Return(nil, assert.AnError).Times(1)

		err := rewardService.RewardUser("test_user_id", 1, 10.0)
		if err == nil {
			t.Errorf("RewardUser() expected error but got nil")
		}
	})
}

func TestRewardServiceImpl_GetRewardHistory(t *testing.T) {
	t.Run("GetRewardHistory", func(t *testing.T) {
		setUpRewardService(t)

		mockedRewardRecordRepository.EXPECT().SearchRewardRecords(
			&repoReal.RewardRecordSearchCondition{
				UserID: "test_user_id",
			},
		).Return([]*model.RewardRecord{
			{
				UserID:        "test_user_id",
				Points:        10.0,
				TaskID:        1,
				OriginPoints:  0.0,
				UpdatedPoints: 10.0,
			},
		}, nil).Times(1)

		rewardRecords, err := rewardService.GetRewardHistory("test_user_id")
		if err != nil {
			t.Errorf("GetRewardHistory() exception = %v", err)
		}

		assert.Equal(t, 1, len(rewardRecords))
		assert.Equal(t, "test_user_id", rewardRecords[0].UserID)
		assert.Equal(t, 10.0, rewardRecords[0].Points)
		assert.Equal(t, 1, rewardRecords[0].TaskID)
	})

	t.Run("GetRewardHistoryFail", func(t *testing.T) {
		setUpRewardService(t)

		mockedRewardRecordRepository.EXPECT().SearchRewardRecords(
			&repoReal.RewardRecordSearchCondition{
				UserID: "test_user_id",
			},
		).Return(nil, assert.AnError).Times(1)

		rewardRecords, err := rewardService.GetRewardHistory("test_user_id")
		if err == nil {
			t.Errorf("GetRewardHistory() expected error but got nil")
		}
		assert.Nil(t, rewardRecords)
	})
}
