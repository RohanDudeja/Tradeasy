package model

import "github.com/golang/protobuf/ptypes/timestamp"

type Holdings struct{
	Userid        string			  `gorm:"foreign_key; column:User" json:"userid"`
	OrderId       string              `gorm:"foreign_key; column:order_id" json:"order_id"`
	Id            int                 `gorm:"primary_key; column:id" json:"id"`
	StockName     string              `gorm:"column:stock_name" json:"stock_name"`
	Quantity      int                 `gorm:"column:quantity" json:"quantity"`
	BuyPrice      int                 `gorm:"column:buy_price" json:"buy_price"`
	OrderedAt     timestamp.Timestamp `gorm:"column:ordered_at" json:"ordered_at" `
	CreatedAt     timestamp.Timestamp `gorm:"column:created_at" json:"created_at" `
	UpdatedAt     timestamp.Timestamp `gorm:"column:updated_at" json:"updated_at" `
	DeletedAt     timestamp.Timestamp `gorm:"column:deleted_at" json:"deleted_at" `
}

func (h *Holdings) TableName() string  {
	return "holdings"
}
