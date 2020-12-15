package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"huaweicloud-sample-rms-sync-resources/models"
)

func QueryMultiResources(provider, resourceType, regionId string, state models.ResourceStateType) ([]models.Resource, error) {
	var resources []models.Resource
	result := DB.Where("provider = ? AND type = ? AND regionId = ? AND state = ?", provider, resourceType, regionId, state).Find(&resources)
	if result.Error != nil {
		return nil, fmt.Errorf("Query resources failed, provider %v, type: %v, regionId: %v",
			provider, resourceType, regionId)
	}
	return resources, nil
}

func deleteOneResource(tx *gorm.DB, resource *models.Resource) error {
	result := tx.Model(resource).
		Where("state = ? AND queryAt < ?", models.ResourceStateNormal, resource.QueryAt).
		Updates(models.Resource{State: models.ResourceStateDeleted, QueryAt: resource.QueryAt})
	if result.Error != nil {
		return result.Error
	}
	logrus.Infof("Delete resource successfully, provider: %v, type: %v, id: %v", resource.Provider, resource.Type, resource.ResourceId)
	return nil
}

func upsertOneResource(tx *gorm.DB, resource *models.Resource) error {
	result := tx.Exec(SQL_TEMPLATE_UPSERT_ONE_RESOURCE, resource.Provider, resource.Type, resource.ResourceId,
		resource.Name, resource.RegionId, resource.ProjectId, resource.ProjectName, resource.EpId, resource.EpName,
		resource.Checksum, resource.Created, resource.Updated, resource.ProvisioningState, resource.Tags,
		resource.Properties, resource.QueryAt, resource.State)
	if result.Error != nil {
		return result.Error
	}
	logrus.Infof("Upsert resource successfully, provider: %v, type: %v, id: %v", resource.Provider, resource.Type, resource.ResourceId)
	return nil
}
