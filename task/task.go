package task

import (
	"github.com/robfig/cron/v3"
	"huaweicloud-sample-rms-sync-resources/config"
	"log"
)

func AddFullSyncJob() {
	cfg := config.GetConfig()
	c := cron.New()
	_, err := c.AddFunc(cfg.Sync.Spec, FullSyncResources)
	if err != nil {
		log.Fatalf("Failed to add full-sync job, err: %v", err)
	}
	c.Start()
}
