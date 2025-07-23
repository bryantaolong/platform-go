package model

import (
	"time"
)

type User struct {
	ID         uint64    `json:"id" gorm:"primaryKey;type:bigserial"`
	Username   string    `json:"username" gorm:"type:varchar(50);unique;not null"`
	Password   string    `json:"-" gorm:"type:varchar(100);not null"`
	Email      string    `json:"email" gorm:"type:varchar(100)"`
	Status     int8      `json:"status" gorm:"type:smallint;default:0"`
	Roles      string    `json:"roles" gorm:"type:varchar(255);default:'ROLE_USER'"`
	CreateTime time.Time `json:"create_time" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time `json:"update_time" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	Deleted    int8      `json:"-" gorm:"type:smallint;default:0;index"`
}

func (User) TableName() string {
	return "user"
}
