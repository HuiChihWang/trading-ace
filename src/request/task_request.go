package request

type GetTaskRequest struct {
	User      string `form:"user_address"`
	Type      string `form:"type"`
	Status    string `form:"status"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
}
