package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"trading-ace/src/database"
	"trading-ace/src/model"
)

var setUpRewardRecordRepo = func(t *testing.T) *rewardRecordRepositoryImpl {
	dbInstance := database.GetDBInstance()

	t.Cleanup(func() {
		dbInstance.Exec("DELETE FROM reward_records")
	})

	return &rewardRecordRepositoryImpl{
		dbInstance: dbInstance,
	}
}

func TestRewardRecordRepositoryImpl_CreateRewardRecord(t *testing.T) {
	t.Run("CreateRewardRecord", func(t *testing.T) {
		repo := setUpRewardRecordRepo(t)
		record := &model.RewardRecord{
			UserID:        "test_user_id",
			Points:        100,
			TaskID:        1,
			OriginPoints:  0,
			UpdatedPoints: 100,
			CreatedAt:     time.Now().UTC(),
		}
		record, err := repo.CreateRewardRecord(record)
		if err != nil {
			t.Errorf("CreateRewardRecord() error = %v", err)
		}

		assert.Equal(t, "test_user_id", record.UserID)
		assert.Equal(t, 100.0, record.Points)
		assert.Equal(t, 0.0, record.OriginPoints)
		assert.Equal(t, 100.0, record.UpdatedPoints)
		assert.Equal(t, 1, record.TaskID)
		assert.NotEmpty(t, record.ID)
	})
}

func TestNewRewardRecordRepositoryImpl_SearchRewardRecords(t *testing.T) {
	t.Run("SearchRewardRecordsOfUser", func(t *testing.T) {
		repo := setUpRewardRecordRepo(t)
		record := &model.RewardRecord{
			UserID:        "test_user_id",
			Points:        100,
			TaskID:        1,
			OriginPoints:  0,
			UpdatedPoints: 100,
			CreatedAt:     time.Now().UTC(),
		}
		record, err := repo.CreateRewardRecord(record)
		if err != nil {
			t.Errorf("CreateRewardRecord() error = %v", err)
		}

		records, err := repo.SearchRewardRecords(&RewardRecordSearchCondition{UserID: "test_user_id"})
		if err != nil {
			t.Errorf("SearchRewardRecords() error = %v", err)
		}

		assert.Equal(t, 1, len(records))
		assert.Equal(t, "test_user_id", records[0].UserID)
		assert.Equal(t, 100.0, records[0].Points)
		assert.Equal(t, 0.0, records[0].OriginPoints)
		assert.Equal(t, 100.0, records[0].UpdatedPoints)
		assert.Equal(t, 1, records[0].TaskID)
		assert.NotEmpty(t, records[0].ID)
	})

	t.Run("SearchRewardRecordsOfUserWithTimeRange", func(t *testing.T) {
		repo := setUpRewardRecordRepo(t)

		records := []*model.RewardRecord{
			{
				UserID:        "test_user_id",
				Points:        100,
				TaskID:        1,
				OriginPoints:  0,
				UpdatedPoints: 100,
				CreatedAt:     time.Now().UTC().Add(-2 * time.Hour), // 2 hours ago
			},
			{
				UserID:        "test_user_id",
				Points:        200,
				TaskID:        2,
				OriginPoints:  50,
				UpdatedPoints: 150,
				CreatedAt:     time.Now().UTC().Add(-1 * time.Hour), // 1 hour ago
			},
			{
				UserID:        "test_user_id",
				Points:        300,
				TaskID:        3,
				OriginPoints:  100,
				UpdatedPoints: 200,
				CreatedAt:     time.Now().UTC().Add(-30 * time.Minute), // 30 minutes ago
			},
			{
				UserID:        "test_user_id",
				Points:        400,
				TaskID:        4,
				OriginPoints:  150,
				UpdatedPoints: 250,
				CreatedAt:     time.Now().UTC().Add(-10 * time.Hour), // 10 hours ago
			},
			{
				UserID:        "other_user_id",
				Points:        500,
				TaskID:        5,
				OriginPoints:  200,
				UpdatedPoints: 300,
				CreatedAt:     time.Now().UTC().Add(-time.Hour * 24), // 1 day ago
			},
		}

		recordsWithID := make([]*model.RewardRecord, 0)
		for _, record := range records {
			newRecord, err := repo.CreateRewardRecord(record)
			if err != nil {
				t.Errorf("CreateRewardRecord() error = %v", err)
			}
			recordsWithID = append(recordsWithID, newRecord)
		}

		searchCondition := &RewardRecordSearchCondition{
			StartTime: time.Now().UTC().Add(-2 * time.Hour),
			Duration:  3 * time.Hour,
		}

		foundRecords, err := repo.SearchRewardRecords(searchCondition)
		if err != nil {
			t.Errorf("SearchRewardRecords() error = %v", err)
		}

		assert.Equal(t, 2, len(foundRecords))

		fundRecordIDs := make([]int, 0)
		for _, record := range foundRecords {
			fundRecordIDs = append(fundRecordIDs, record.ID)
		}

		assert.Equal(t, []int{recordsWithID[2].ID, recordsWithID[1].ID}, fundRecordIDs)
	})

	t.Run("SearchRecordShouldReturnNilSliceWhenNoQueriedData", func(t *testing.T) {
		repo := setUpRewardRecordRepo(t)
		records, err := repo.SearchRewardRecords(&RewardRecordSearchCondition{UserID: "test_user_id"})
		if err != nil {
			t.Errorf("SearchRewardRecords() error = %v", err)
		}

		assert.Empty(t, records)
	})
}
