package model

import (
	"github.com/jinzhu/gorm"
	"time"
)
// Stocks ... model for stocks
type Stocks struct {
	gorm.Model
	StockTickerSymbol     string    `gorm:"primary_key" json:"stock_ticker_symbol,omitempty"`
	StockName             string    `json:"stock_name,omitempty"`
	LTP                   int       `json:"ltp,omitempty"`
	OpenPrice             int       `json:"open_price,omitempty"`
	HighPrice             int       `json:"high_price,omitempty"`
	LowPrice              int       `json:"low_price,omitempty"`
	PreviousDayClose      int       `json:"previous_day_close,omitempty"`
	PercentageChange      int       `json:"percentage_change,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	DeletedAt             time.Time `json:"deleted_at"`
}

// BuyOrderBook ... book that stores buy orders
type BuyOrderBook struct {
	gorm.Model
	OrderID           string    `json:"order_id,omitempty" gorm:"primary_key"`
	StockTickerSymbol string    `json:"stock_ticker_symbol,omitempty"`
	OrderQuantity     int       `json:"order_quantity,omitempty"`
	OrderStatus       string    `json:"order_status,omitempty"`
	OrderPrice        int       `json:"order_price,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         time.Time `json:"deleted_at"`
}

// SellOrderBook ... book that stores sell orders
type SellOrderBook struct {
	gorm.Model
	OrderID           string    `json:"order_id,omitempty" gorm:"primary_key"`
	StockTickerSymbol string    `json:"stock_ticker_symbol,omitempty"`
	OrderQuantity     int       `json:"order_quantity,omitempty"`
	OrderStatus       string    `json:"order_status,omitempty"`
	OrderPrice        int       `json:"order_price,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         time.Time `json:"deleted_at"`
}

func (s *Stocks) TableName() string{
	return "stocks"
}

func (b *BuyOrderBook) TableName() string{
	return "buy_order_book"
}

func (s *SellOrderBook) TableName() string{
	return "sell_order_book"
}