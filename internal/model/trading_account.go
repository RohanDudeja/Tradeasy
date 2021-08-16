package model

import (
	"time"
)

type TradingAccount struct {
	UserId       string    `json:"user_id,omitempty" gorm:"column:user_id"`
	Id           int       `gorm:"column:id;primary_key" json:"id,omitempty"`
	PanCardNo    string    `gorm:"column:pan_card_no" json:"pan_card_no,omitempty"`
	BankAccNo    string    `gorm:"column:bank_acc_no" json:"bank_acc_no,omitempty"`
	TradingAccId string    `gorm:"column:trading_acc_no" json:"trading_acc_id,omitempty"`
	Balance      int64     `gorm:"column:balance" json:"balance,omitempty"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    time.Time `json:"deleted_at" gorm:"column:deleted_at;default:NULL"`
}

func (u *TradingAccount) TableName() string {
	return "trading_account"
}
