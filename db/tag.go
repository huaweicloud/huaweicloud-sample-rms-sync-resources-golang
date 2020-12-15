package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"huaweicloud-sample-rms-sync-resources/models"
)

func queryMultiTags(tx *gorm.DB, provider, resourceType, resourceId string) ([]models.Tag, error) {
	var tags []models.Tag
	result := tx.Where("provider = ? AND type = ? AND resourceId = ?", provider, resourceType, resourceId).Find(&tags)
	if result.Error != nil {
		return nil, fmt.Errorf("Query tags failed, provider %v, type: %v, resourceId: %v", provider, resourceType, resourceId)
	}
	return tags, nil
}

func insertOneTag(tx *gorm.DB, tag *models.Tag) error {
	result := tx.Create(tag)
	if result.Error != nil {
		return result.Error
	}
	logrus.Infof("Insert tag successfully, provider: %v, type: %v, resourceId: %v, name: %v", tag.Provider, tag.Type, tag.ResourceId, tag.Name)
	return nil
}

func updateOneTag(tx *gorm.DB, tag *models.Tag) error {
	result := tx.Save(tag)
	if result.Error != nil {
		return result.Error
	}
	logrus.Infof("Update tag successfully, provider: %v, type: %v, resourceId: %v, name: %v", tag.Provider, tag.Type, tag.ResourceId, tag.Name)
	return nil
}

func deleteOneTag(tx *gorm.DB, tag *models.Tag) error {
	result := tx.Delete(tag)
	if result.Error != nil {
		return result.Error
	}
	logrus.Infof("Delete tag successfully, provider: %v, type: %v, resourceId: %v, name: %v", tag.Provider, tag.Type, tag.ResourceId, tag.Name)
	return nil
}
