package main

import (
	"github.com/sirupsen/logrus"
	"huaweicloud-sample-rms-sync-resources/core"
	"huaweicloud-sample-rms-sync-resources/task"
	"log"
	"net/http"
)

func main() {
	task.AddFullSyncJob()
	logrus.Info("Server started..")
	// This interface is only used in test scenarios and can be manually invoked to trigger full synchronization
	http.HandleFunc("/v1/sync", core.HandleFullSync)
	// This interface is used to accept SMN messages, and verify the signature of SMN messages. Others have no right to call
	http.HandleFunc("/v1/smn/notify", core.HandleSmnMessage)
	log.Fatal(http.ListenAndServe(":7777", nil))
}
