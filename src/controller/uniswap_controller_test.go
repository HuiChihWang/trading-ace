package controller

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"trading-ace/mock/service"
	"trading-ace/src/contract"
)

type uniSwapEventControllerTestSuite struct {
	uniSwapController    UniSwapEventController
	mockedUniSwapService *service.MockUniSwapService
}

func (s *uniSwapEventControllerTestSuite) setUp(t *testing.T) {
	s.mockedUniSwapService = service.NewMockUniSwapService(t)
	s.uniSwapController = &uniSwapEventController{
		uniSwapService: s.mockedUniSwapService,
	}
}

func TestHandleUniSwapV2Event(t *testing.T) {
	testSuite := &uniSwapEventControllerTestSuite{}

	testSender := "0x0000000000000000000000000000001234567890"
	testReiciver := "0x00000000000000000000000000000056767890"

	t.Run("HandleUniSwapV2Event - USDC to WETH", func(t *testing.T) {
		testSuite.setUp(t)
		testEvent := &contract.UniSwapV2SwapEvent{
			Amount0In:  big.NewInt(0),
			Amount0Out: big.NewInt(123456),
			Amount1In:  big.NewInt(1234567890),
			Amount1Out: big.NewInt(0),
			Sender:     common.HexToAddress(testSender),
			To:         common.HexToAddress(testReiciver),
		}

		testSuite.mockedUniSwapService.EXPECT().ProcessUniSwapTransaction(
			testSender,
			0.123456).
			Return(nil).Times(1)
		err := testSuite.uniSwapController.HandleUniSwapV2Event(testEvent)
		assert.Nil(t, err)
	})

	t.Run("HandleUniSwapV2Event - WETH to USDC", func(t *testing.T) {
		testSuite.setUp(t)

		testEvent := &contract.UniSwapV2SwapEvent{
			Amount0In:  big.NewInt(123456),
			Amount0Out: big.NewInt(0),
			Amount1In:  big.NewInt(0),
			Amount1Out: big.NewInt(1234567890),
			Sender:     common.HexToAddress(testSender),
			To:         common.HexToAddress(testReiciver),
		}

		testSuite.mockedUniSwapService.EXPECT().ProcessUniSwapTransaction(testSender, 0.123456).Return(nil).Times(1)
		err := testSuite.uniSwapController.HandleUniSwapV2Event(testEvent)
		assert.Nil(t, err)
	})

	t.Run("HandleUniSwapV2Event - nil event", func(t *testing.T) {
		testSuite.setUp(t)
		err := testSuite.uniSwapController.HandleUniSwapV2Event(nil)
		assert.Nil(t, err)
	})

	t.Run("HandleUniSwapV2Event - error", func(t *testing.T) {
		testSuite.setUp(t)
		testEvent := &contract.UniSwapV2SwapEvent{
			Amount0In:  big.NewInt(123456),
			Amount0Out: big.NewInt(0),
			Amount1In:  big.NewInt(0),
			Amount1Out: big.NewInt(1234567890),
			Sender:     common.HexToAddress(testSender),
			To:         common.HexToAddress(testReiciver),
		}

		testSuite.mockedUniSwapService.EXPECT().ProcessUniSwapTransaction(testSender, 0.123456).Return(assert.AnError).Times(1)
		err := testSuite.uniSwapController.HandleUniSwapV2Event(testEvent)
		assert.NotNil(t, err)
	})
}
