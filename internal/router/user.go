package model

import (
	"github.com/jinzhu/gorm"
	"math/big"
	"time"
)

type Users struct {
	Userid    string    `gorm:"column:user_id;primary_key" json:"userid,omitempty"`
	EmailId   string    `json:"emailId,omitempty" gorm:"column:email_id"`
	Password  string    `json:"password,omitempty" gorm:"column:password"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (u *Users) TableName() string {
	return "users"
}
