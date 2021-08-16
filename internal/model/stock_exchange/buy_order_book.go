package model

import (
	"time"
)

// BuyOrderBook ... schema for order book that stores buy orders
type BuyOrderBook struct {
	ID                int       `json:"id" gorm:"primary_key; column:id"`
	OrderID           string    `json:"order_id,omitempty" gorm:"column:order_id"`
	StockTickerSymbol string    `json:"stock_ticker_symbol,omitempty" gorm:"column:stock_ticker_symbol"`
	OrderQuantity     int       `json:"order_quantity,omitempty" gorm:"column:order_quantity"`
	OrderStatus       string    `json:"order_status,omitempty" gorm:"column:order_status"`
	OrderPrice        int       `json:"order_price,omitempty" gorm:"column:order_price"`
	OrderType         string    `json:"order_type,omitempty" gorm:"column:order_type"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt         time.Time `json:"deleted_at" gorm:"column:deleted_at; default:NULL"`
}

func (b *BuyOrderBook) TableName() string {
	return "buy_order_book"
}
