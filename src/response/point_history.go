package response

import (
	"trading-ace/src/model"
)

type PointHistory struct {
	User              string  `json:"user_address"`
	DistributedPoints float64 `json:"distributed_points"`
	TotalPoints       float64 `json:"total_points"`
	UpdatedAt         string  `json:"updated_at"`
}

type PointHistoryCollection []*PointHistory

func CreatePointHistory(record *model.RewardRecord) *PointHistory {
	return &PointHistory{
		User:              record.UserID,
		DistributedPoints: record.Points,
		TotalPoints:       record.UpdatedPoints,
		UpdatedAt:         record.CreatedAt.String(),
	}
}

func CreatePointHistoryCollection(recordList *[]*model.RewardRecord) *PointHistoryCollection {
	var collection PointHistoryCollection
	for _, record := range *recordList {
		collection = append(collection, CreatePointHistory(record))
	}

	if collection == nil {
		collection = make(PointHistoryCollection, 0)
	}

	return &collection
}
