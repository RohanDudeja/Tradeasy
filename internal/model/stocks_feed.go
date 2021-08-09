package model

import (
	"time"
)

type StocksFeed struct {
	Id        int64   `json:"id" gorm:"primary_key;column:id"`
	StockName string    `json:"stock_name" gorm:"column:stock_name"`
	LTP       int       `json:"ltp" gorm:"column:ltp"`
	Open      int       `json:"open" gorm:"column:open"`
	High      int       `json:"high" gorm:"column:high"`
	Low       int       `json:"low" gorm:"column:low"`
	TradedAt  time.Time `json:"traded_at" gorm:"column:traded_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"column:deleted_at;default:NULL"`
}

func (s *StocksFeed) TableName() string {
	return "stocks_feed"
}
