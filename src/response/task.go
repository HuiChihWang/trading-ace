package response

import (
	"time"
	"trading-ace/src/model"
)

type TaskCollection []*Task

type Task struct {
	User              string    `json:"user_address"`
	Type              string    `json:"type"`
	Status            string    `json:"status"`
	SwapAmount        float64   `json:"swap_amount"`
	DistributedPoints float64   `json:"distributed_points"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewTask(task *model.Task, distributedPoints float64) *Task {
	return &Task{
		User:              task.UserID,
		Type:              string(task.Type),
		Status:            string(task.Status),
		SwapAmount:        task.SwapAmount,
		DistributedPoints: distributedPoints,
		CreatedAt:         task.CreatedAt,
	}
}
