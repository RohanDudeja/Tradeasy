package Model

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"go/types"
)

type PendingOrders struct {
	Userid string
	User Users `gorm:"foreignKey:Userid"`
	OrderId string `gorm:"primaryKey"`
	StockName string `gorm:"type:varchar"`
	OrderType string `gorm:"type:varchar"`
	BookType string `gorm:"type:varchar"`
	LimitPrice int
	Quantity int
	OrderPrice int
	Status string `gorm:"type:varchar"`
	CreatedAt timestamp.Timestamp
	UpdatedAt timestamp.Timestamp
	DeletedAt timestamp.Timestamp
}

type Holdings struct{
	Userid string
	User Users `gorm:"foreignKey:Userid"`
	OrderId string
	PendingOrders PendingOrders `gorm:"foreignKey:OrderId"`
	Id int `gorm:"primaryKey"`
	StockName string `gorm:"type:varchar"`
	Quantity int
	BuyPrice int
	OrderedAt timestamp.Timestamp
	CreatedAt timestamp.Timestamp
	UpdatedAt timestamp.Timestamp
	DeletedAt timestamp.Timestamp
}

type OrderHistory struct {
	Userid string
	User Users `gorm:"foreignKey:Userid"`
	OrderId string
	PendingOrders PendingOrders `gorm:"foreignKey:OrderId"`
	Id int `gorm:"primaryKey"`
	StockName string `gorm:"type:varchar"`
	Quantity int
	BuyPrice int
	SellPrice int
	CommissionFee int
	BoughtAt timestamp.Timestamp
	SoldAt timestamp.Timestamp
	CreatedAt timestamp.Timestamp
	UpdatedAt timestamp.Timestamp
	DeletedAt timestamp.Timestamp
}

func (p *PendingOrders) TableName() string  {
	return "PendingOrders"
}

func (h *Holdings) TableName() string  {
	return "Holdings"
}

func (o *OrderHistory) TableName() string  {
	return "OrderHistory"
}

