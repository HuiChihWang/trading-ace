package model

import "time"

type RewardRecord struct {
	ID            int       `json:"id"`
	UserID        string    `json:"user_id"`
	Points        float64   `json:"points"`
	TaskID        int       `json:"task_id"`
	OriginPoints  float64   `json:"origin_points"`
	UpdatedPoints float64   `json:"updated_points"`
	CreatedAt     time.Time `json:"created_at"`
}
