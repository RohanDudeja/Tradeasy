package model

import (
	"time"
)

type Holdings struct {
	UserId    string    `json:"user_id" gorm:"column:user_id"`
	OrderId   string    `json:"order_id" gorm:"column:order_id"`
	Id        int       `json:"id" gorm:"primary_key; column:id"`
	StockName string    `json:"stock_name" gorm:"column:stock_name"`
	Quantity  int       `json:"quantity" gorm:"column:quantity"`
	BuyPrice  int       `json:"buy_price" gorm:"column:buy_price"`
	OrderedAt time.Time `json:"ordered_at" gorm:"column:ordered_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"column:deleted_at;default:NULL"`
}

func (h *Holdings) TableName() string {
	return "holdings"
}
