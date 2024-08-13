package contract

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"strings"
)

type UniSwapV2SwapEvent struct {
	Sender     common.Address
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	To         common.Address
}

type UniSwapV2Contract struct {
	contractAddress common.Address
	abi             *abi.ABI
	client          *ethclient.Client
}

func NewUniSwapV2Contract(contractAddressStr string, abiPath string, client *ethclient.Client) (*UniSwapV2Contract, error) {
	contractAddress := common.HexToAddress(contractAddressStr)

	abiBytes, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return nil, err
	}

	return &UniSwapV2Contract{
		contractAddress: contractAddress,
		abi:             &parsedABI,
		client:          client,
	}, nil
}

func (c *UniSwapV2Contract) ListenSwapEvents(callback func(event *UniSwapV2SwapEvent) error) {
	eventChan := make(chan *UniSwapV2SwapEvent)

	go func() {
		query := ethereum.FilterQuery{
			Addresses: []common.Address{c.contractAddress},
			Topics:    [][]common.Hash{{c.abi.Events["Swap"].ID}},
		}

		logs := make(chan types.Log)
		sub, err := c.client.SubscribeFilterLogs(context.Background(), query, logs)
		if err != nil {
			log.Fatalf("Failed to subscribe to filter logs: %v", err)
		}

		for {
			select {
			case err := <-sub.Err():
				log.Fatalf("Subscription exception: %v", err)
			case vLog := <-logs:
				var event UniSwapV2SwapEvent

				err := c.abi.UnpackIntoInterface(&event, "Swap", vLog.Data)
				if err != nil {
					log.Printf("Failed to unpack event data: %v", err)
					continue
				}

				if len(vLog.Topics) != 3 {
					log.Printf("Invalid number of topics: %v", len(vLog.Topics))
					continue
				}

				event.Sender = common.HexToAddress(vLog.Topics[1].Hex())
				event.To = common.HexToAddress(vLog.Topics[2].Hex())

				eventChan <- &event
			}
		}
	}()

	go func() {
		for event := range eventChan {
			err := callback(event)
			if err != nil {
				log.Printf("Callback exception: %v", err)
			}
		}
	}()
}
