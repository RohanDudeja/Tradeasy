package model

import (
	"time"
)

type PendingOrders struct {
	UserId     string    `gorm:"column:user_id" json:"user_id" `
	OrderId    string    `gorm:"primary_key;column:order_id" json:"order_id"`
	StockName  string    `gorm:"column:stock_name" json:"stock_name"`
	OrderType  string    `gorm:"column:order_type" json:"order_type"`
	BookType   string    `gorm:"column:book_type" json:"book_type"`
	LimitPrice int       `gorm:"column:limit_price" json:"limit_price" `
	Quantity   int       `gorm:"column:quantity" json:"quantity" `
	OrderPrice int       `gorm:"column:order_price" json:"order_price" `
	Status     string    `gorm:"column:status" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at" `
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at" `
	DeletedAt  time.Time `gorm:"column:deleted_at; default:NULL" json:"deleted_at" `
}

func (p *PendingOrders) TableName() string {
	return "pending_orders"
}
