package models

type Tag struct {
	Provider   string   `gorm:"primaryKey"`
	Type       string   `gorm:"primaryKey"`
	ResourceId string   `gorm:"primaryKey;column:resourceId"`
	Name       string   `gorm:"type:VARCHAR(20);primaryKey;index:idx_name_value,priority:1;not null"`
	Value      string   `gorm:"type:VARCHAR(20);index:idx_name_value,priority:2;not null"`
	Resource   Resource `gorm:"foreignKey:Provider,Type,ResourceId"`
}
