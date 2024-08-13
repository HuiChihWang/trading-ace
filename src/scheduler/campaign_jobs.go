package scheduler

import (
	"github.com/go-co-op/gocron/v2"
	"time"
)

type CampaignCallback func(start time.Time, end time.Time) error

func CreateCampaignJobs(s gocron.Scheduler, startTime time.Time, weeks int, callback CampaignCallback) {
	weekDuration := time.Hour * 24 * 7
	for i := 0; i < weeks; i++ {
		start := startTime.Add(weekDuration * time.Duration(i))
		end := start.Add(weekDuration)

		_, _ = s.NewJob(
			gocron.OneTimeJob(
				gocron.OneTimeJobStartDateTime(end),
			),
			gocron.NewTask(callback, start, end),
			gocron.WithName("campaign job"),
		)
	}
}
