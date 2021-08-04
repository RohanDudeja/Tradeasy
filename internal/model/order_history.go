package model

import (
	"time"
)

type OrderHistory struct {
	UserId        string              `gorm:"foreign_key; column:user_id" json:"user_id" `
	OrderId       string              `gorm:"foreign_key; column:order_id" json:"order_id"`
	Id            int                 `gorm:"primaryKey; column:id" json:"id"`
	StockName     string              `gorm:"column:stock_name" json:"stock_name"`
	Quantity      int                 `gorm:"column:quantity" json:"quantity" `
	BuyPrice      int                 `gorm:"column:buy_price" json:"buy_price" `
	SellPrice     int                 `gorm:"column:sell_price" json:"sell_price" `
	CommissionFee int                 `gorm:"column:commission_fee" json:"commission_fee" `
	BoughtAt      time.Time `gorm:"column:bought_at" json:"bought_at" `
	SoldAt        time.Time `gorm:"column:sold_at" json:"sold_at" `
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at" `
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at" `
	DeletedAt     time.Time `gorm:"column:deleted_at" json:"deleted_at" `
}

func (o *OrderHistory) TableName() string {
	return "order_history"
}
