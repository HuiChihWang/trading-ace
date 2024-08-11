package model

import "time"

type RewardRecord struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Points    float64   `json:"points"`
	TaskID    int       `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRewardRecord(userID string, points float64, taskID int) *RewardRecord {
	return &RewardRecord{
		UserID:    userID,
		Points:    points,
		TaskID:    taskID,
		CreatedAt: time.Now(),
	}
}
