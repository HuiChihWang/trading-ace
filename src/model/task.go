package model

import (
	"database/sql"
	"time"
)

type TaskStatus string

const (
	TaskStatusPending TaskStatus = "pending"
	TaskStatusDone    TaskStatus = "done"
)

type TaskType string

const (
	TaskTypeOnboarding TaskType = "on_boarding"
	TaskTypeSharedPool TaskType = "shared_pool"
)

type Task struct {
	ID          int          `json:"id"`
	Status      TaskStatus   `json:"status"`
	Type        TaskType     `json:"type"`
	UserID      string       `json:"user_id"`
	SwapAmount  float64      `json:"swap_amount"`
	CreatedAt   time.Time    `json:"created_at"`
	CompletedAt sql.NullTime `json:"completed_at"`
}

func NewTask(userID string, taskType TaskType, swapAmount float64) *Task {
	return &Task{
		Status:      TaskStatusPending,
		UserID:      userID,
		Type:        taskType,
		SwapAmount:  swapAmount,
		CreatedAt:   time.Now(),
		CompletedAt: sql.NullTime{},
	}
}
