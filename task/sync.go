package task

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"huaweicloud-sample-rms-sync-resources/client"
	"huaweicloud-sample-rms-sync-resources/config"
	"huaweicloud-sample-rms-sync-resources/db"
	"huaweicloud-sample-rms-sync-resources/helper"
	"huaweicloud-sample-rms-sync-resources/models"
	"strings"
)

func FullSyncResources() {
	logrus.Infof("Start to full sync resources..")
	cfg := config.GetConfig()
	service := cfg.Rms

	resourceTypeToRegions, err := loadResourceTypeToRegionMap()
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, rt := range service.ResourceTypes {
		for _, regionId := range *resourceTypeToRegions[rt] {
			go syncSpecificTypeResourcesInRegion(rt, regionId)
		}
	}
}

func loadResourceTypeToRegionMap() (map[string]*[]string, error) {
	resourceTypeToRegions := make(map[string]*[]string)
	resp, err := client.ListProviders()
	if err != nil {
		return nil, fmt.Errorf("Failed to list providers from RMS, err: %v", err)
	}
	providers := resp.ResourceProviders
	for _, provider := range *providers {
		for _, resourceType := range *provider.ResourceTypes {
			joinedTags := []string{*provider.Provider, *resourceType.Name}
			resourceTypeToRegions[strings.Join(joinedTags, ".")] = resourceType.Regions
		}
	}
	return resourceTypeToRegions, nil
}

func syncSpecificTypeResourcesInRegion(resourceType, regionId string) {
	logrus.Infof("Start to sync resources, resourceType: %v, regionId: %v", resourceType, regionId)
	contents := strings.Split(resourceType, ".")
	provider, resourceType := contents[0], contents[1]
	resourceMapFromRM, err1 := helper.ListResourcesFromRM(provider, resourceType, regionId)
	if err1 != nil {
		logrus.Errorf("Failed to list resource from RM, err: %v", err1)
		return
	}
	resourceMapFromDB, err2 := helper.ListNormalResourcesFromDB(provider, resourceType, regionId, models.ResourceStateNormal)
	if err2 != nil {
		logrus.Errorf("Failed to list resources from DB, err: %v", err2)
		return
	}

	resourcesToBeCreated, resourcesToBeUpdated, resourcesToBeDeleted := distinguishResources(resourceMapFromRM, resourceMapFromDB)
	logrus.Infof("%d resources will be created, resourceType: %v, regionId: %v", len(resourcesToBeCreated), resourceType, regionId)
	logrus.Infof("%d resources will be updated, resourceType: %v, regionId: %v", len(resourcesToBeUpdated), resourceType, regionId)
	logrus.Infof("%d resources will be deleted, resourceType: %v, regionId: %v", len(resourcesToBeDeleted), resourceType, regionId)
	go dealWithResourcesToBeCreated(resourcesToBeCreated)
	go dealWithResourceToBeUpdated(resourcesToBeUpdated)
	go dealWithResourcesToBeDeleted(resourcesToBeDeleted)
}

func distinguishResources(resourceMapFromRM map[string]*models.Resource, resourceMapFromDB map[string]*models.Resource) ([]*models.Resource, []*models.Resource, []*models.Resource) {
	var resourcesToBeCreated []*models.Resource
	var resourcesToBeUpdated []*models.Resource
	var resourcesToBeDeleted []*models.Resource
	for resourceIdFromRM, resourceFromRM := range resourceMapFromRM {
		var resourceFromDB *models.Resource
		if resourceMapFromDB[resourceIdFromRM] == nil {
			resourcesToBeCreated = append(resourcesToBeCreated, resourceFromRM)
			delete(resourceMapFromRM, resourceIdFromRM)
		} else {
			resourceFromDB = resourceMapFromDB[resourceIdFromRM]
			if resourceFromDB.Checksum != resourceFromRM.Checksum && resourceFromDB.QueryAt.Before(resourceFromRM.QueryAt) {
				resourcesToBeUpdated = append(resourcesToBeUpdated, resourceFromRM)
			}
			delete(resourceMapFromRM, resourceIdFromRM)
			delete(resourceMapFromDB, resourceIdFromRM)
		}
	}
	for _, v := range resourceMapFromDB {
		resourcesToBeDeleted = append(resourcesToBeDeleted, v)
	}
	return resourcesToBeCreated, resourcesToBeUpdated, resourcesToBeDeleted
}

func dealWithResourcesToBeCreated(resourcesTobeCreated []*models.Resource) {
	for _, resource := range resourcesTobeCreated {
		if err := db.CreateOrUpdateOneResource(resource); err != nil {
			logrus.Error(err)
		}
	}
}

func dealWithResourceToBeUpdated(resourcesToBeUpdated []*models.Resource) {
	for _, resource := range resourcesToBeUpdated {
		if err := db.CreateOrUpdateOneResource(resource); err != nil {
			logrus.Error(err)
		}
	}
}

func dealWithResourcesToBeDeleted(resourcesTobeDeleted []*models.Resource) {
	for _, resource := range resourcesTobeDeleted {
		if err := db.DeleteOneResource(resource); err != nil {
			logrus.Error(err)
		}
	}
}
