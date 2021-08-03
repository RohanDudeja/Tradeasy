package Model

import "time"

type Payments struct {
	Id             string    `gorm:"primary_key; column:id" json:"id"`
	Userid         string    `gorm:"foreign_key; column:user_id" json:"user_id"`
	RazorpayLinkId string    `gorm:"column:razorpay_link_id" json:"razorpay_link_id"`
	RazorpayLink   string    `gorm:"column:razorpay_link" json:"razorpay_link"`
	Amount         int       `gorm:"column:amount" json:"amount"`
	Type           string    `gorm:"column:type" json:"type"`
	CurrentAmount  int       `gorm:"column:current_amount" json:"current_amount"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (b *Payments) TableName() string {
	return "payments"
}
