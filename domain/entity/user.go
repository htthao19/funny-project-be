package entity

import (
	"time"
)

// User model.
type User struct {
	ID    uint   `gorm:"primary_key;column:id;auto_increment:true"`
	Name  string `gorm:"type:varchar(100);column:name"`
	Email string `gorm:"type:varchar(100);column:email;index:unique"`

	CreatedAt time.Time `gorm:"column:created_at;autocreatetime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoupdatetime"`
}

// TableName is the pluralized version of struct name
func (User) TableName() string {
	return "user"
}
