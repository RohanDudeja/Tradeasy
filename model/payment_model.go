package model

import (
"github.com/golang/protobuf/ptypes/timestamp"
)

type Payments struct {
	id   string  `gorm:"primaryKey"`
	Userid string
	User Users `gorm:"foreignKey:Userid"`
	razorpayUserId string `gorm:"type:varchar"`
	razorpayLinkId string  `gorm:"type:varchar"`
	Amount int
	Type string `gorm:"type:varchar"`
	currentAmount int
	CreatedAt timestamp.Timestamp
	UpdatedAt timestamp.Timestamp
	DeletedAt timestamp.Timestamp
}

func (b *Payments) TableName() string {
	return "payments"
}