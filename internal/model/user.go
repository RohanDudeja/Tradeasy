package model

import (
	"time"
)

type Users struct {
	UserId    string    `gorm:"column:user_id;primary_key" json:"user_id,omitempty"`
	EmailId   string    `json:"email_id,omitempty" gorm:"column:email_id"`
	Password  string    `json:"password,omitempty" gorm:"column:password"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"column:deleted_at;default:NULL"`
}

func (u *Users) TableName() string {
	return "users"
}
