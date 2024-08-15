package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"trading-ace/mock/service"
	"trading-ace/src/model"
	"trading-ace/src/response"
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

		testUser := "test_user_id"
		startTimeStr := "2024-08-15T00:00:00Z"
		endTimeStr := "2024-08-16T23:59:59Z"
		testUrl := fmt.Sprintf("/api/reward-history?user_address=%s&start_time=%s&end_time=%s", testUser, startTimeStr, endTimeStr)

		testResponseWriter := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(testResponseWriter)
		testContext.Request = httptest.NewRequest(http.MethodGet, testUrl, nil)

		startTime, _ := time.Parse(time.RFC3339, startTimeStr)
		endTime, _ := time.Parse(time.RFC3339, endTimeStr)

		testSuite.mockedRewardService.EXPECT().
			GetRewardHistory(testUser, startTime, endTime.Sub(startTime)).
			Return(recordHistory, nil)

		testSuite.rewardController.GetRewardHistoryOfUser(testContext)

		assert.Equal(t, http.StatusOK, testContext.Writer.Status())

		expectedRes := response.CreatePointHistoryCollection(&recordHistory)

		expectedResStr, err := json.Marshal(expectedRes)

		assert.Nil(t, err)
		assert.Equal(t, string(expectedResStr), testResponseWriter.Body.String())
	})

	t.Run("GetRewardHistoryOfUser with error", func(t *testing.T) {
		testSuite.setUp(t)

		testUser := "test_user_id"
		startTimeStr := "2024-08-15T00:00:00Z"
		endTimeStr := "2024-08-16T23:59:59Z"
		testUrl := fmt.Sprintf("/api/reward-history?user_address=%s&start_time=%s&end_time=%s", testUser, startTimeStr, endTimeStr)

		testResponseWriter := httptest.NewRecorder()
		testContext, _ := gin.CreateTestContext(testResponseWriter)
		testContext.Request = httptest.NewRequest(http.MethodGet, testUrl, nil)

		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		endTime, err := time.Parse(time.RFC3339, endTimeStr)

		if err != nil {
			t.Errorf("Error parsing time: %v", err)
		}

		testSuite.mockedRewardService.EXPECT().
			GetRewardHistory("test_user_id", startTime, endTime.Sub(startTime)).
			Return(nil, assert.AnError)

		testSuite.rewardController.GetRewardHistoryOfUser(testContext)

		assert.Equal(t, http.StatusInternalServerError, testContext.Writer.Status())

		var exception map[string]string
		err = json.Unmarshal(testResponseWriter.Body.Bytes(), &exception)

		assert.Nil(t, err)
		assert.Equal(t, assert.AnError.Error(), exception["exception"])
	})
}
