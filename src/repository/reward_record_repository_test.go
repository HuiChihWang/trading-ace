package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"trading-ace/src/config"
	"trading-ace/src/database"
)

func TestRewardRecordRepositoryImpl(t *testing.T) {
	setUpRewardRecordRepo := func(t *testing.T) *rewardRecordRepositoryImpl {
		dbInstance := database.CreateDBInstance(&config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5435",
			Username: "postgres",
			Password: "postgres",
			DBName:   "trading_ace_test",
		})

		t.Cleanup(func() {
			dbInstance.Exec("DELETE FROM reward_records")
			fmt.Println("[Tear Down] Cleaned up reward_records table")
		})

		return &rewardRecordRepositoryImpl{
			dbInstance: dbInstance,
		}
	}

	t.Run("CreateRewardRecord", func(t *testing.T) {
		repo := setUpRewardRecordRepo(t)
		record, err := repo.CreateRewardRecord("test_user_id", 100, 1)
		if err != nil {
			t.Errorf("CreateRewardRecord() error = %v", err)
		}

		assert.Equal(t, "test_user_id", record.UserID)
		assert.Equal(t, 100.0, record.Points)
		assert.Equal(t, 1, record.TaskID)
		assert.Equal(t, record.CreatedAt.Location(), time.UTC)
		assert.NotEmpty(t, record.ID)
	})

	t.Run("GetRewardRecordsByUserID", func(t *testing.T) {
		repo := setUpRewardRecordRepo(t)
		for i := 0; i < 10; i++ {
			_, _ = repo.CreateRewardRecord("test_user_id", 100, 1)
		}

		records, err := repo.GetRewardRecordsByUserID("test_user_id")

		if err != nil {
			t.Errorf("GetRewardRecordsByUserID() error = %v", err)
		}

		assert.Equal(t, 10, len(records))
		for _, record := range records {
			assert.Equal(t, "test_user_id", record.UserID)
			assert.Equal(t, 100.0, record.Points)
			assert.Equal(t, 1, record.TaskID)
			assert.NotEmpty(t, record.CreatedAt)
			assert.NotEmpty(t, record.ID)
		}
	})
}
