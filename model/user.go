package model

import (
	"github.com/jinzhu/gorm"
	"math/big"
	"time"
)
type Users struct {
	gorm.Model
	Userid           string    `gorm:"primary_key" json:"userid,omitempty"`
	EmailId          string    `json:"emailId,omitempty"`
	Password              string   `json:"password,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

func (u *Users) TableName() string{
	return "users"
}


