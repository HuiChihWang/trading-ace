package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"trading-ace/mock/service"
	"trading-ace/src/model"
)

type rewardControllerTestSuite struct {
	rewardController    RewardController
	mockedRewardService *service.MockRewardService
}

func (s *rewardControllerTestSuite) setUp(t *testing.T) {
	s.mockedRewardService = service.NewMockRewardService(t)
	s.rewardController = &rewardController{
		rewardService: s.mockedRewardService,
	}
}

func TestGetRewardHistoryOfUser(t *testing.T) {
	testSuite := &rewardControllerTestSuite{}
	recordHistory := []*model.RewardRecord{
		{
			ID:     1,
			UserID: "test_user_id",
			Points: 10.0,
			TaskID: 1,
		},
		{
			ID:     2,
			UserID: "test_user_id",
			Points: 10.0,
			TaskID: 1,
		},
		{
			ID:     3,
			UserID: "test_user_id",
			Points: 10.0,
			TaskID: 1,
		},
		{
			ID:     4,
			UserID: "test_user_id",
			Points: 10.0,
			TaskID: 1,
		},
	}

	t.Run("GetRewardHistoryOfUser", func(t *testing.T) {
		testSuite.setUp(t)

		testSuite.mockedRewardService.EXPECT().
			GetRewardHistory("test_user_id").
			Return(recordHistory, nil)

		testResponseWriter := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(testResponseWriter)
		testContext.Params = gin.Params{
			{Key: "userId", Value: "test_user_id"},
		}

		testSuite.rewardController.GetRewardHistoryOfUser(testContext)

		assert.Equal(t, http.StatusOK, testContext.Writer.Status())

		var recordHistoryFromRes []*model.RewardRecord
		err := json.Unmarshal(testResponseWriter.Body.Bytes(), &recordHistoryFromRes)

		assert.Nil(t, err)
		assert.Equal(t, recordHistory, recordHistoryFromRes)
	})

	t.Run("GetRewardHistoryOfUser with error", func(t *testing.T) {
		testSuite.setUp(t)

		testSuite.mockedRewardService.EXPECT().
			GetRewardHistory("test_user_id").
			Return(nil, assert.AnError)

		testResponseWriter := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(testResponseWriter)
		testContext.Params = gin.Params{
			{Key: "userId", Value: "test_user_id"},
		}

		testSuite.rewardController.GetRewardHistoryOfUser(testContext)

		assert.Equal(t, http.StatusInternalServerError, testContext.Writer.Status())

		var exception map[string]string
		err := json.Unmarshal(testResponseWriter.Body.Bytes(), &exception)

		assert.Nil(t, err)
		assert.Equal(t, assert.AnError.Error(), exception["exception"])
	})
}
