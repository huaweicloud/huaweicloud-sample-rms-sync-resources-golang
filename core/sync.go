package core

import (
	"huaweicloud-sample-rms-sync-resources/task"
	"net/http"
)

func HandleFullSync(w http.ResponseWriter, r *http.Request) {
	go task.FullSyncResources()
}
