package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type UserWatchlist struct {
	gorm.Model
	ID          int    `gorm:"primary_key" json:"id,omitempty"`
	Userid      string `gorm:"foreign_key:Userid" json:"userid,omitempty"`
	WatchListId string `gorm:"watchlist_id" json:"watchlist_id,omitempty"`
	StockName   string `gorm:"stock_name" json:"stock_name,omitempty"`
}

func (u *UserWatchlist) TableName() string {
	return "user_watchlist"
}
