package job

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

type Type string

const (
	TypeUniSwapTransaction Type = "uni_swap:transaction"
)

type UniSwapTransactionPayload struct {
	SenderID   string  `json:"sender_id"`
	SwapAmount float64 `json:"swap_amount"`
}

func NewUniSwapTransactionTask(payload *UniSwapTransactionPayload) (*asynq.Task, error) {
	return createAsyncQTask(TypeUniSwapTransaction, payload)
}

func createAsyncQTask(jobType Type, payload interface{}) (*asynq.Task, error) {
	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(string(jobType), payloadByte), nil
}
