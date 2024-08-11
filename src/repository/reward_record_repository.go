package repository

import (
	"database/sql"
	"fmt"
	"time"
	"trading-ace/src/database"
	"trading-ace/src/model"
)

const rewardRecordTableName = "reward_records"

type RewardRecordRepository interface {
	CreateRewardRecord(userID string, points float64, taskID int) (*model.RewardRecord, error)
	GetRewardRecordsByUserID(userID string) ([]*model.RewardRecord, error)
}

type rewardRecordRepositoryImpl struct {
	dbInstance *sql.DB
}

func NewRewardRecordRepository() RewardRecordRepository {
	a := database.GetDBInstance()
	return &rewardRecordRepositoryImpl{
		dbInstance: a,
	}
}

func (r *rewardRecordRepositoryImpl) CreateRewardRecord(userID string, points float64, taskID int) (*model.RewardRecord, error) {
	record := model.NewRewardRecord(userID, points, taskID)
	record.CreatedAt = record.CreatedAt.In(time.UTC)

	sqlCommand := fmt.Sprintf("INSERT INTO %s (user_id, points, task_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id", rewardRecordTableName)
	stmt, err := r.dbInstance.Prepare(sqlCommand)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(record.UserID, record.Points, record.TaskID, record.CreatedAt).Scan(&record.ID)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r *rewardRecordRepositoryImpl) GetRewardRecordsByUserID(userID string) ([]*model.RewardRecord, error) {
	sqlCommand := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", rewardRecordTableName)
	rows, err := r.dbInstance.Query(sqlCommand, userID)

	if err != nil {
		return nil, err
	}

	var records []*model.RewardRecord
	for rows.Next() {
		var record model.RewardRecord
		err := rows.Scan(&record.ID, &record.UserID, &record.Points, &record.TaskID, &record.CreatedAt)
		if err != nil {
			return nil, err
		}

		records = append(records, &record)
	}

	return records, nil

}
