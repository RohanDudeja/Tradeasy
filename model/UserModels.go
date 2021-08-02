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

type UsersTradingAccDetails struct {
	gorm.Model
	Userid       string    `json:"userid,omitempty" gorm:"foreign_key:Userid"`
	Id           int       `json:"id,omitempty"`
	PanCardNo    string    `json:"panCardNo,omitempty"`
	BankAccNo    string    `json:"bank_acc_no,omitempty"`
	TradingAccId string       `json:"trading_acc_id,omitempty"`
	Balance      big.Int `json:"balance,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

func (u *Users) TableName() string{
	return "users"
}

func (u *UsersTradingAccDetails) TableName() string{
	return "users_trading_acc_details"
}
