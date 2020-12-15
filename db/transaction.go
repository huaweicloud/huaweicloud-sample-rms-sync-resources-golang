package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"huaweicloud-sample-rms-sync-resources/models"
)

func CreateOrUpdateOneResource(resource *models.Resource) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := upsertOneResource(tx, resource); err != nil {
			return fmt.Errorf("Failed to upsert a resource, provider: %v, type: %v, resourceId: %v",
				resource.Provider, resource.Type, resource.ResourceId)
		}
		if err := updateTags(tx, resource); err != nil {
			return fmt.Errorf("Failed to update tags, provider: %v, type: %v, resourceId: %v",
				resource.Provider, resource.Type, resource.ResourceId)
		}
		logrus.Infof("Transaction CreateOrUpdateOneResource finished, provider: %v, type: %v, resourceId: %v",
			resource.Provider, resource.Type, resource.ResourceId)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Transaction CreateOrUpdateOneResource failed, provider: %v, type: %v, resourceId: %v, err: %v",
			resource.Provider, resource.Type, resource.ResourceId, err)
	}
	return nil
}

func DeleteOneResource(resource *models.Resource) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := deleteOneResource(tx, resource); err != nil {
			return fmt.Errorf("Failed to delele a resource, provider: %v, type: %v, resourceId: %v", resource.Provider, resource.Type, resource.ResourceId)
		}
		if err := deleteTags(tx, resource); err != nil {
			return fmt.Errorf("Failed to delete tags, provider: %v, type: %v, resourceId: %v", resource.Provider, resource.Type, resource.ResourceId)
		}
		logrus.Infof("Transaction DeleteOneResource finished, provider: %v, type: %v, resourceId: %v",
			resource.Provider, resource.Type, resource.ResourceId)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Transaction DeleteOneResource failed, provider: %v, type: %v, resourceId: %v, err: %v",
			resource.Provider, resource.Type, resource.ResourceId, err)
	}
	return nil
}

func deleteTags(tx *gorm.DB, resource *models.Resource) error {
	tagsFromDB, err := queryMultiTags(tx, resource.Provider, resource.Type, resource.ResourceId)
	if err != nil {
		return err
	}
	for i := 0; i < len(tagsFromDB); i++ {
		if err := deleteOneTag(tx, &tagsFromDB[i]); err != nil {
			return err
		}
	}
	return nil
}

func updateTags(tx *gorm.DB, resource *models.Resource) error {
	tagsFromDB, err := queryMultiTags(tx, resource.Provider, resource.Type, resource.ResourceId)
	if err != nil {
		return err
	}
	tagsToBeCreated := []*models.Tag{}
	tagsToBeUpdated := []*models.Tag{}
	tagMapFromRM := resource.Tags
	tagMapFromDB := make(map[string]string)
	for _, tag := range tagsFromDB {
		tagMapFromDB[tag.Name] = tag.Value
	}

	for name, value := range tagMapFromRM {
		tag := models.Tag{Name: name, Value: value, Provider: resource.Provider, Type: resource.Type, ResourceId: resource.ResourceId}
		if tagMapFromDB[name] == "" {
			tagsToBeCreated = append(tagsToBeCreated, &tag)
		} else {
			tagsToBeUpdated = append(tagsToBeUpdated, &tag)
			delete(tagMapFromDB, name)
		}
		delete(tagMapFromRM, name)
	}

	for name, value := range tagMapFromDB {
		if err := deleteOneTag(tx, &models.Tag{Name: name, Value: value, Provider: resource.Provider, Type: resource.Type, ResourceId: resource.ResourceId}); err != nil {
			return err
		}
	}
	for _, tag := range tagsToBeCreated {
		if err := insertOneTag(tx, tag); err != nil {
			return err
		}
	}
	for _, tag := range tagsToBeUpdated {
		if err := updateOneTag(tx, tag); err != nil {
			return err
		}
	}
	return nil
}
