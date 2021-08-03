package model

import "github.com/golang/protobuf/ptypes/timestamp"

type OrderHistory struct {
	Userid        string              `gorm:"foreign_key; column:userid" json:"userid" `
	OrderId       string              `gorm:"foreign_key; column:order_id" json:"order_id"`
	Id            int                 `gorm:"primaryKey; column:id" json:"id"`
	StockName     string              `gorm:"column:stock_name" json:"stock_name"`
	Quantity      int                 `gorm:"column:quantity" json:"quantity" `
	BuyPrice      int                 `gorm:"column:buy_price" json:"buy_price" `
	SellPrice     int                 `gorm:"column:sell_price" json:"sell_price" `
	CommissionFee int                 `gorm:"column:commission_fee" json:"commission_fee" `
	BoughtAt      timestamp.Timestamp `gorm:"column:bought_at" json:"bought_at" `
	SoldAt        timestamp.Timestamp `gorm:"column:sold_at" json:"sold_at" `
	CreatedAt     timestamp.Timestamp `gorm:"column:created_at" json:"created_at" `
	UpdatedAt     timestamp.Timestamp `gorm:"column:updated_at" json:"updated_at" `
	DeletedAt     timestamp.Timestamp `gorm:"column:deleted_at" json:"deleted_at" `
}

func (o *OrderHistory) TableName() string {
	return "order_history"
}
