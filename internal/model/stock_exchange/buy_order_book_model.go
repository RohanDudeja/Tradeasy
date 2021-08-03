package stock_exchange

import (
	"gorm.io/gorm"
)

// BuyOrderBook ... schema for order book that stores buy orders
type BuyOrderBook struct {
	gorm.Model
	OrderID           string    `json:"order_id,omitempty" gorm:"primary_key; column:order_id"`
	StockTickerSymbol string    `json:"stock_ticker_symbol,omitempty" gorm:"column:stock_ticker_symbol"`
	OrderQuantity     int       `json:"order_quantity,omitempty" gorm:"column:order_quantity"`
	OrderStatus       string    `json:"order_status,omitempty" gorm:"column:order_status"`
	OrderPrice        int       `json:"order_price,omitempty" gorm:"column:order_price"`
}

func (b *BuyOrderBook) TableName() string{
	return "buy_order_book"
}