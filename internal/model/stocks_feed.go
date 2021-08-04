package model

import (
	"math/big"
	"time"
)

type StocksFeed struct {
	Id        big.Int             `gorm:"primary_key;column:id" json:"id"`
	StockName string              `gorm:"column:stock_name" json:"stock_name" `
	LTP       int                 `gorm:"column:ltp" json:"ltp" `
	Open      int                 `gorm:"column:open" json:"open" `
	High      int                 `gorm:"column:high" json:"high" `
	Low       int                 `gorm:"column:low" json:"low" `
	TradedAt  time.Time `gorm:"column:traded_at" json:"traded_at" `
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" `
	DeletedAt time.Time `gorm:"column:deleted-at" json:"deleted_at" `
}

func (s *StocksFeed) TableName() string {
	return "stocks_feed"
}
