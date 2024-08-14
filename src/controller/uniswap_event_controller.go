package controller

import (
	"fmt"
	"sync"
	"trading-ace/src/contract"
	"trading-ace/src/service"
)

type UniSwapEventController interface {
	HandleUniSwapV2Event(event *contract.UniSwapV2SwapEvent) error
}

type uniSwapEventController struct {
	uniSwapService service.UniSwapService
}

var (
	uniSwapEventControllerInstance *uniSwapEventController
	uniSwapEventControllerOnce     sync.Once
)

func GetUniSwapEventControllerInstance() UniSwapEventController {
	uniSwapEventControllerOnce.Do(func() {
		uniSwapEventControllerInstance = &uniSwapEventController{
			uniSwapService: service.NewUniSwapService(),
		}
	})
	return uniSwapEventControllerInstance
}

func (u *uniSwapEventController) HandleUniSwapV2Event(event *contract.UniSwapV2SwapEvent) error {
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

	return u.uniSwapService.ProcessUniSwapTransaction(senderID, swapAmountFloat)
}
