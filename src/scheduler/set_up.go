package scheduler

import (
	"github.com/go-co-op/gocron/v2"
	"log"
	"trading-ace/src/config"
	"trading-ace/src/service"
)

func SetUpScheduler() (gocron.Scheduler, error) {
	sch, err := gocron.NewScheduler()

	if err != nil {
		return nil, err
	}

	campaignConfig := config.GetAppConfig().Campaign

	if campaignConfig != nil {
		CreateCampaignJobs(sch, campaignConfig.GetCampaignStartTime(), campaignConfig.Weeks, service.NewUniSwapService().ProcessSharedPool)
	}

	sch.Start()

	return sch, nil
}

func ShutDowScheduler(sch gocron.Scheduler) {
	if sch == nil {
		return
	}

	err := sch.Shutdown()
	if err != nil {
		log.Fatalf("failed to shutdown scheduler")
	}
}
