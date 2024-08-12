package controller

import (
	"errors"
	"fmt"
	"log"
	"trading-ace/src/contract"
	"trading-ace/src/exception"
	"trading-ace/src/model"
	"trading-ace/src/service"
)

func HandleUniSwapV2Event(event *contract.UniSwapV2SwapEvent) error {
	userService := service.NewUserService()

	userID := event.Sender.String()
	var user *model.User
	user, err := userService.GetUserByID(userID)

	if err != nil && !errors.Is(err, exception.UserNotFoundError) {
		return err
	}

	if user == nil {
		user, err = userService.CreateUser(userID)
	}

	if err != nil {
		return err
	}

	usdcAmountIn := event.Amount0In.Int64()
	usdcAmountOut := event.Amount0Out.Int64()

	var swapAmount int64
	if usdcAmountIn == 0 {
		swapAmount = usdcAmountOut
	} else {
		swapAmount = usdcAmountIn
	}

	swapAmountFloat := float64(swapAmount) / 1e6

	fmt.Printf("Swap Event:\n")
	fmt.Printf("Sender: %s\n", event.Sender.Hex())
	fmt.Printf("Swap Amount: %f USD\n", swapAmountFloat)

	taskService := service.NewTaskService()

	if !taskService.IsUserOnboardingCompleted(userID) {
		err := taskService.ProcessOnBoarding(userID, swapAmountFloat)

		if err != nil {
			log.Println("Onboarding process failed" + err.Error())
		}
	} else {
		_, err := taskService.CreateTask(userID, model.TaskTypeSharedPool, swapAmountFloat)

		if err != nil {
			log.Println("Shared pool task creation failed" + err.Error())
		}
	}

	return nil
}
