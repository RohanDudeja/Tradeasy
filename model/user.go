package model

import (
	"github.com/jinzhu/gorm"
	"math/big"
	"time"
)

type Users struct {
	gorm.Model
	Userid   string `gorm:"primary_key" json:"userid,omitempty"`
	EmailId  string `json:"emailId,omitempty"`
	Password string `json:"password,omitempty"`
}

func (u *Users) TableName() string {
	return "users"
}
