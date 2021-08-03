package Model

import "github.com/golang/protobuf/ptypes/timestamp"

type Payments struct {
	id             string              `gorm:"primary_key; column:id" json:"id"`
	Userid         string              `gorm:"foreign_key; column:user_id" json:"user_id"`
	razorpayUserId string              `gorm:"column:razorpay_user_id" json:"razorpay_user_id"`
	razorpayLinkId string              `gorm:"column:razorpay_user_id" json:"razorpay_link_id"`
	Amount         int                 `gorm:"column:amount" json:"amount"`
	Type           string              `gorm:"column:type" json:"type"`
	currentAmount  int                 `gorm:"column:current_amount" json:"current_amount"`
	CreatedAt      timestamp.Timestamp `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      timestamp.Timestamp `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      timestamp.Timestamp `gorm:"column:deleted_at" json:"deleted_at"`
}

func (b *Payments) TableName() string {
	return "payments"
}