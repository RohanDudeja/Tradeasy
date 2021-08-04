package model

import (
	"github.com/jinzhu/gorm"
)
// Stocks ... schema for stocks data in stock exchange
type Stocks struct {
	ID                    uint      `gorm:"primary_key;column:id" json:"id"`
	StockTickerSymbol     string    `gorm:"column:stock_ticker_symbol" json:"stock_ticker_symbol,omitempty"`
	StockName             string    `json:"stock_name,omitempty" gorm:"column:stock_name"`
	LTP                   int       `json:"ltp,omitempty" gorm:"column:ltp"`
	OpenPrice             int       `json:"open_price,omitempty" gorm:"column:open_price"`
	HighPrice             int       `json:"high_price,omitempty" gorm:"column:high_price"`
	LowPrice              int       `json:"low_price,omitempty" gorm:"column:low_price"`
	PreviousDayClose      int       `json:"previous_day_close,omitempty" gorm:"column:previous_day_close"`
	PercentageChange      int       `json:"percentage_change,omitempty" gorm:"column:percentage_change"`
	CreatedAt             time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt             time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (s *Stocks) TableName() string{
	return "stocks"
}
