package request

type GetRewordHistoryRequest struct {
	User      string `form:"user_address"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
}
