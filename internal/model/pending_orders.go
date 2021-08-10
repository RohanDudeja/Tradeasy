package model

import (
	"time"
)

type PendingOrders struct {
	UserId     string    `json:"user_id" gorm:"column:user_id"`
	OrderId    string    `json:"order_id" gorm:"primary_key;column:order_id"`
	StockName  string    `json:"stock_name" gorm:"column:stock_name"`
	OrderType  string    `json:"order_type" gorm:"column:order_type"`
	BookType   string    `json:"book_type" gorm:"column:book_type"`
	LimitPrice int       `json:"limit_price" gorm:"column:limit_price"`
	Quantity   int       `json:"quantity" gorm:"column:quantity"`
	OrderPrice int       `json:"order_price" gorm:"column:order_price"`
	Status     string    `json:"status" gorm:"column:status"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt  time.Time `json:"deleted_at" gorm:"column:deleted_at;default:NULL"`
}

func (p *PendingOrders) TableName() string {
	return "pending_orders"
}
