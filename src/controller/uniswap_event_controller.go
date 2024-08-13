package controller

import (
	"fmt"
	"trading-ace/src/contract"
	"trading-ace/src/service"
)

func HandleUniSwapV2Event(event *contract.UniSwapV2SwapEvent) error {
	if event == nil {
		return nil
	}

	usdcAmountIn := event.Amount0In.Int64()
	usdcAmountOut := event.Amount0Out.Int64()

	var swapAmount int64
	if usdcAmountIn == 0 {
		swapAmount = usdcAmountOut
	} else {
		swapAmount = usdcAmountIn
	}

	senderID := event.Sender.String()
	swapAmountFloat := float64(swapAmount) / 1e6

	fmt.Printf("Swap Event:\n")
	fmt.Printf("Sender: %s\n", senderID)
	fmt.Printf("Swap Amount: %f USD\n", swapAmountFloat)

	return service.NewUniSwapService().ProcessUniSwapTransaction(senderID, swapAmountFloat)
}
