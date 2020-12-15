package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rms/v1/model"
	"time"
)

type Tags map[string]string
type Properties map[string]interface{}

func (tags Tags) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Failed to scan tags, value is not []byte, value: %v", value)
	}
	return json.Unmarshal(b, &tags)
}

func (tags Tags) Value() (driver.Value, error) {
	if tags == nil {
		return nil, nil
	}
	return json.Marshal(tags)
}

func (properties Properties) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Failed to scan properties, value is not []byte, value: %v", value)
	}
	return json.Unmarshal(b, &properties)
}

func (properties Properties) Value() (driver.Value, error) {
	if properties == nil {
		return nil, nil
	}
	return json.Marshal(properties)
}

type Resource struct {
	Provider          string            `gorm:"type:VARCHAR(20);primaryKey;index:idx_provider_type_regionId_state,priority:1"`
	Type              string            `gorm:"type:VARCHAR(20);primaryKey;index:idx_provider_type_regionId_state,priority:2"`
	ResourceId        string            `gorm:"type:VARCHAR(64);primaryKey;column:resourceId" json:"id" mapstructure:"id"`
	Name              string            `gorm:"type:VARCHAR(64);not null"`
	RegionId          string            `gorm:"type:VARCHAR(20);not null;index:idx_provider_type_regionId_state,priority:3;column:regionId" json:"region_id" mapstructure:"region_id"`
	ProjectId         string            `gorm:"type:VARCHAR(64);not null;column:projectId" json:"project_id" mapstructure:"project_id"`
	ProjectName       string            `gorm:"type:VARCHAR(64);not null;column:projectName" json:"project_name" mapstructure:"project_name"`
	EpId              string            `gorm:"type:VARCHAR(64);not null;column:epId" json:"ep_id" mapstructure:"ep_id"`
	EpName            string            `gorm:"type:VARCHAR(64);not null;column:epName" json:"ep_name" mapstructure:"ep_name"`
	Checksum          string            `gorm:"type:VARCHAR(64);not null"`
	Created           string            `gorm:"type:VARCHAR(30)"`
	Updated           string            `gorm:"type:VARCHAR(30)"`
	ProvisioningState string            `gorm:"type:VARCHAR(10);column:provisioningState" json:"provisioning_state" mapstructure:"provisioning_state"`
	Tags              Tags              `gorm:"type:TEXT"`
	Properties        Properties        `gorm:"type:TEXT"`
	QueryAt           time.Time         `gorm:"type:DATETIME(3);not null;column:queryAt"`
	State             ResourceStateType `gorm:"type:VARCHAR(10);not null;index:idx_provider_type_regionId_state,priority:4"`
}

func ConvertResourceEntityToResource(resourceEntity *model.ResourceEntity, queryAt *time.Time) *Resource {
	resource := new(Resource)
	resource.Provider = *resourceEntity.Provider
	resource.Type = *resourceEntity.Type
	resource.ResourceId = *resourceEntity.Id
	resource.Name = *resourceEntity.Name
	resource.RegionId = *resourceEntity.RegionId
	resource.ProjectId = *resourceEntity.ProjectId
	resource.ProjectName = *resourceEntity.ProjectName
	resource.EpId = *resourceEntity.EpId
	resource.EpName = *resourceEntity.EpName
	resource.Checksum = *resourceEntity.Checksum
	resource.Created = *resourceEntity.Created
	resource.Updated = *resourceEntity.Updated
	resource.ProvisioningState = *resourceEntity.ProvisioningState
	resource.Tags = resourceEntity.Tags
	resource.Properties = resourceEntity.Properties
	resource.QueryAt = *queryAt
	resource.State = ResourceStateNormal
	return resource
}
