package Transaction_engine

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"math/big"
)

type StocksFeed struct {
	Id big.Int `gorm:"primary_key;column:id" json:"id""`
	StockName string `json:"stock_name"`
	LTP int `json:"ltp" gorm:"column:ltp"`
	Open int `json:"open" gorm:"column:open"`
	High int `json:"high" gorm:"column:high"`
	Low int `json:"low" gorm:"column:low"`
	TradedAt timestamp.Timestamp `json:"traded_at" gorm:"column:traded_at"`
	CreatedAt timestamp.Timestamp `json:"created_at" gorm:"column:created_at"`
	DeletedAt timestamp.Timestamp `json:"deleted_at" gorm:"column:deleted-at"`
}

func (s *StocksFeed) TableName() string  {
	return "StocksFeed"
}
