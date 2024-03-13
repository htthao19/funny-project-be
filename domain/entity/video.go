package entity

import (
	"time"
)

// Video model.
type Video struct {
	ID          uint   `gorm:"primary_key;column:id;auto_increment:true"`
	URL         string `gorm:"type:varchar(200);column:url"`
	SharedBy    string `gorm:"type:varchar(100);column:shared_by"`
	Description string `gorm:"column:description"`

	CreatedAt time.Time `gorm:"column:created_at;autocreatetime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoupdatetime"`
}

// TableName is the pluralized version of struct name
func (Video) TableName() string {
	return "video"
}
