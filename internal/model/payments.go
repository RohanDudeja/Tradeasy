package model

import (
	"time"
)

type Payments struct {
	Id             int       `json:"id" gorm:"primary_key; column:id"`
	UserId         string    `json:"user_id" gorm:"column:user_id"`
	RazorpayLinkId string    `json:"razorpay_link_id" gorm:"column:razorpay_link_id"`
	RazorpayLink   string    `json:"razorpay_link" gorm:"column:razorpay_link"`
	Amount         int64     `json:"amount" gorm:"column:amount"`
	PaymentType    string    `json:"payment_type" gorm:"column:payment_type"`
	CurrentBalance int64     `json:"current_balance" gorm:"column:current_balance"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt      time.Time `json:"deleted_at" gorm:"column:deleted_at;default:NULL"`
}

func (b *Payments) TableName() string {
	return "payments"
}
