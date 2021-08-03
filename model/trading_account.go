package model

import (
	"github.com/jinzhu/gorm"
	"math/big"
	"time"
)

type TradingAccount struct {
	gorm.Model
	Userid       string  `json:"userid,omitempty" gorm:"foreign_key:Userid"`
	Id           int     `gorm:"primary_key" json:"id,omitempty"`
	PanCardNo    string  `gorm:"pan_card_no" json:"panCardNo,omitempty"`
	BankAccNo    string  `gorm:"bank_acc_no" json:"bank_acc_no,omitempty"`
	TradingAccId string  `gorm:"trading_acc_no" json:"trading_acc_id,omitempty"`
	Balance      big.Int `gorm:"balance" json:"balance,omitempty"`
}

func (u *TradingAccount) TableName() string {
	return "trading_account"
}
