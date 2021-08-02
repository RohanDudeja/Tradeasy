package model

import "gorm.io/gorm"

// BuyOrder for stock exchange engine
type BuyOrder struct {
	gorm.Model
	OrderID string `json:"order_id"`
	StockTickerSymbol string `json:"stock_ticker_symbol"`
	Quantity uint `json:"quantity"`
	Status string `json:"status"`
	Price uint `json:"price"`
}

// SellOrder for stock exchange engine
type SellOrder struct {
	gorm.Model
	OrderID string `json:"order_id"`
	StockTickerSymbol string `json:"stock_ticker_symbol"`
	Quantity uint `json:"quantity"`
	Status string `json:"status"`
	Price uint `json:"price"`
}

func (b *BuyOrder) TableName() string {
	return  "Buy Orders"
}
func (s *SellOrder) TableName() string{
	return "Sell Orders"
}