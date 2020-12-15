package models

type Notification struct {
	NotificationType         string      `json:"notification_type"`
	NotificationCreationTime string      `json:"notification_creation_time"`
	DomainId                 string      `json:"domain_id"`
	Detail                   interface{} `json:"detail"`
}

type ResourceChangedDetail struct {
	ResourceId   string    `mapstructure:"resource_id"`
	ResourceType string    `mapstructure:"resource_type"`
	EventType    EventType `mapstructure:"event_type"`
	CaptureTime  string    `mapstructure:"capture_time"`
	Resource     *Resource `mapstructure:"resource"`
}

type ResourceRelationChangedDetail struct {
	ResourceId   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	EventType    EventType `json:"event_type"`
	CaptureTime  string    `json:"capture_time"`
}

type SnapshotArchiveCompletedNotification struct {
	SnapshotId string    `json:"snapshot_id"`
	RegionId   string    `json:"region_id"`
	BucketName string    `json:"bucket_name"`
	ObjectKeys *[]string `json:"object_keys"`
}

type NotificationArchiveCompletedNotification struct {
	RegionId   string `json:"region_id"`
	BucketName string `json:"bucket_name"`
	ObjectKey  string `json:"object_key"`
}
