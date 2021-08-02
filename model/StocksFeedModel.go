package Model

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"math/big"
)

type StocksFeed struct {
	Id big.Int `gorm:"primaryKey"`
	StockName string `gorm:"type:varchar"`
	LTP int
	Open int
	High int
	Low int
	TradedAt timestamp.Timestamp
	CreatedAt timestamp.Timestamp
	DeletedAt timestamp.Timestamp
}

func (s *StocksFeed) TableName() string  {
	return "StocksFeed"
}
