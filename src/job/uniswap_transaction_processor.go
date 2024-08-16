package job

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
	"trading-ace/src/service"
)

type UniSwapTransactionProcessor struct {
	uniSwapService service.UniSwapService
}

func NewUniSwapTransactionProcessor() *UniSwapTransactionProcessor {
	return &UniSwapTransactionProcessor{
		uniSwapService: service.NewUniSwapService(),
	}
}

func (processor *UniSwapTransactionProcessor) ProcessTask(_ context.Context, t *asynq.Task) error {
	var payload UniSwapTransactionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	senderID := payload.SenderID
	swapAmount := payload.SwapAmount

	log.Println("Processing UniSwap transaction for senderID: ", senderID, " swapAmount: ", swapAmount)

	return processor.uniSwapService.ProcessUniSwapTransaction(senderID, swapAmount)
}
