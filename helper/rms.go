package helper

import (
	rms_model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rms/v1/model"
	"huaweicloud-sample-rms-sync-resources/client"
	"huaweicloud-sample-rms-sync-resources/config"
	"huaweicloud-sample-rms-sync-resources/db"
	"huaweicloud-sample-rms-sync-resources/models"
	"time"
)

var limit int32

func init() {
	cfg := config.GetConfig()
	limit = cfg.Rms.Limit
}

func ListResourcesFromRM(provider, resourceType, regionId string) (map[string]*models.Resource, error) {
	resourceMap := make(map[string]*models.Resource)
	var response *rms_model.ListResourcesResponse
	var err error
	queryAt := time.Now().UTC()
	response, err = client.ListLimitedResources(provider, resourceType, regionId, limit, nil)
	if err != nil {
		return nil, err
	}
	for _, resourceEntity := range *response.Resources {
		resource := models.ConvertResourceEntityToResource(&resourceEntity, &queryAt)
		resourceMap[resource.ResourceId] = resource
	}
	for response.PageInfo.NextMarker != nil {
		queryAt := time.Now().UTC()
		response, err = client.ListLimitedResources(provider, resourceType, regionId, limit, response.PageInfo.NextMarker)
		if err != nil {
			return nil, err
		}
		for _, resourceEntity := range *response.Resources {
			resource := models.ConvertResourceEntityToResource(&resourceEntity, &queryAt)
			resourceMap[resource.ResourceId] = resource
		}
	}
	return resourceMap, nil
}

func ListNormalResourcesFromDB(provider, resourceType, regionId string, state models.ResourceStateType) (map[string]*models.Resource, error) {
	resourceMap := make(map[string]*models.Resource)
	resources, err := db.QueryMultiResources(provider, resourceType, regionId, state)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(resources); i++ {
		resourceMap[resources[i].ResourceId] = &resources[i]
	}
	return resourceMap, nil
}
