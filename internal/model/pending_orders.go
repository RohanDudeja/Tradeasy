package model

import "github.com/golang/protobuf/ptypes/timestamp"

type PendingOrders struct {
	Userid     string              `gorm:"foreign_key; column:userid" json:"userid" `
	OrderId    string              `gorm:"primary_key;column:order_id" json:"order_id"`
	StockName  string              `gorm:"column:stock_name" json:"stock_name"`
	OrderType  string              `gorm:"column:order_type" json:"order_type"`
	BookType   string              `gorm:"column:book_type" json:"book_type"`
	LimitPrice int                 `gorm:"column:limit_price" json:"limit_price" `
	Quantity   int                 `gorm:"column:quantity" json:"quantity" `
	OrderPrice int                 `gorm:"column:order_price" json:"order_price" `
	Status     string              `gorm:"column:status" json:"status"`
	CreatedAt  timestamp.Timestamp `gorm:"column:created_at" json:"created_at" `
	UpdatedAt  timestamp.Timestamp `gorm:"column:updated_at" json:"updated_at" `
	DeletedAt  timestamp.Timestamp `gorm:"column:deleted_at" json:"deleted_at" `
}

func (p *PendingOrders) TableName() string {
	return "pending_orders"
}
