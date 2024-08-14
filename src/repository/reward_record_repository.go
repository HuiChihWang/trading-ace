package repository

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"time"
	"trading-ace/src/database"
	"trading-ace/src/model"
)

const rewardRecordTableName = "reward_records"

type RewardRecordSearchCondition struct {
	StartTime time.Time
	Duration  time.Duration
	UserID    string
}

type RewardRecordRepository interface {
	CreateRewardRecord(rewardRecord *model.RewardRecord) (*model.RewardRecord, error)
	SearchRewardRecords(condition *RewardRecordSearchCondition) ([]*model.RewardRecord, error)
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

func (r *rewardRecordRepositoryImpl) CreateRewardRecord(rewardRecord *model.RewardRecord) (*model.RewardRecord, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlCommand, args, err := psql.Insert(rewardRecordTableName).
		Columns("user_id", "points", "task_id", "created_at", "original_points", "updated_points").
		Values(rewardRecord.UserID, rewardRecord.Points, rewardRecord.TaskID, rewardRecord.CreatedAt, rewardRecord.OriginPoints, rewardRecord.UpdatedPoints).
		Suffix("RETURNING id").ToSql()

	fmt.Printf("SQL Command: %s\n", sqlCommand)
	fmt.Printf("Arguments: %v\n", args)

	err = r.dbInstance.QueryRow(sqlCommand, args...).Scan(&rewardRecord.ID)

	if err != nil {
		return nil, err
	}

	return rewardRecord, nil
}

func (r *rewardRecordRepositoryImpl) SearchRewardRecords(condition *RewardRecordSearchCondition) ([]*model.RewardRecord, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	query := psql.
		Select("id, user_id, points, task_id, created_at, original_points, updated_points").
		From(rewardRecordTableName)

	if condition.UserID != "" {
		query = query.Where(squirrel.Eq{"user_id": condition.UserID})
	}

	if !condition.StartTime.IsZero() && condition.Duration != 0 {
		query = query.Where(squirrel.Gt{"created_at": condition.StartTime})
		query = query.Where(squirrel.Lt{"created_at": condition.StartTime.Add(condition.Duration)})
	}

	sqlCommand, args, err := query.OrderBy("id DESC").ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.dbInstance.Query(sqlCommand, args...)
	if err != nil {
		return nil, err
	}

	var records []*model.RewardRecord
	for rows.Next() {
		var record model.RewardRecord
		err := rows.Scan(&record.ID, &record.UserID, &record.Points, &record.TaskID, &record.CreatedAt, &record.OriginPoints, &record.UpdatedPoints)
		if err != nil {
			return nil, err
		}

		records = append(records, &record)
	}

	return records, nil
}
