package model

import (
	"github.com/jinzhu/gorm"
	"math/big"
	"time"
)

type TradingAccount struct {
	gorm.Model
	Userid       string    `json:"userid,omitempty" gorm:"foreign_key:Userid"`
	Id           int       `gorm:"primary_key" json:"id,omitempty"`
	PanCardNo    string    `json:"panCardNo,omitempty"`
	BankAccNo    string    `json:"bank_acc_no,omitempty"`
	TradingAccId string    `json:"trading_acc_id,omitempty"`
	Balance      big.Int   `json:"balance,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

func (u *TradingAccount) TableName() string {
	return "trading_account"
}
