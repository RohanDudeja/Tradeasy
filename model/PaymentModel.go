package model

import (
"github.com/golang/protobuf/ptypes/timestamp"
)

type Payments struct {
	id   string  `json:"id"`
	Userid  string  `json:"user_id"`
	Razorpay_user_id string  `json:"razorpay_user_id"`
	Razorpay_link_id string   `json:"razorpay_link_id"`
	Amount float64   `json:"amount"`
	Type string       `json:"type"`
	CurrentAmount string `json:"current_amount"`
	Created_at timestamp.Timestamp  `json:"created_at"`
	Updated_at timestamp.Timestamp   `json:"updated_at"`
	Deleted_at timestamp.Timestamp    `json:"deleted_at"`
}

func (b *Payments) TableName() string {
	return "payments"
}