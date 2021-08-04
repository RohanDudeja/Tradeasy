package model

import (
	"time"
)

type Holdings struct{
	UserId        string			  `gorm:"foreign_key; column:user_id" json:"user_id"`
	OrderId       string              `gorm:"foreign_key; column:order_id" json:"order_id"`
	Id            int                 `gorm:"primary_key; column:id" json:"id"`
	StockName     string              `gorm:"column:stock_name" json:"stock_name"`
	Quantity      int                 `gorm:"column:quantity" json:"quantity"`
	BuyPrice      int                 `gorm:"column:buy_price" json:"buy_price"`
	OrderedAt     time.Time `gorm:"column:ordered_at" json:"ordered_at" `
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at" `
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at" `
	DeletedAt     time.Time `gorm:"column:deleted_at" json:"deleted_at" `
}

func (h *Holdings) TableName() string  {
	return "holdings"
}
