package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

// SellOrderBook ... schema for order book that stores sell orders
type SellOrderBook struct {
	gorm.Model
	OrderID           string    `json:"order_id,omitempty" gorm:"primary_key; column:order_id"`
	StockTickerSymbol string    `json:"stock_ticker_symbol,omitempty" gorm:"column:stock_ticker_symbol"`
	OrderQuantity     int       `json:"order_quantity,omitempty" gorm:"column:order_quantity"`
	OrderStatus       string    `json:"order_status,omitempty" gorm:"column:order_status"`
	OrderPrice        int       `json:"order_price,omitempty" gorm:"column:order_price"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt         time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (s *SellOrderBook) TableName() string{
	return "sell_order_book"
}
