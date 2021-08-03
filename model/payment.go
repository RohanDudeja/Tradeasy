package Model

import "github.com/golang/protobuf/ptypes/timestamp"

type Payments struct {
	id             string              `gorm:"primaryKey" json:"id"`
	Userid         string              `json:"userid"`
	User           Users               `gorm:"foreignKey:Userid" json:"user"`
	razorpayUserId string              `gorm:"type:varchar" json:"razorpay_user_id"`
	razorpayLinkId string              `gorm:"type:varchar" json:"razorpay_link_id"`
	Amount         int                 `json:"amount"`
	Type           string              `gorm:"type:varchar" json:"type"`
	currentAmount  int                 `json:"current_amount"`
	CreatedAt      timestamp.Timestamp `json:"created_at"`
	UpdatedAt      timestamp.Timestamp `json:"updated_at"`
	DeletedAt      timestamp.Timestamp `json:"deleted_at"`
}

func (b *Payments) TableName() string {
	return "payments"
}
