package model

import (
	"time"
)

type UserWatchlist struct {
	Id          int       `gorm:"column:id;primary_key" json:"id,omitempty"`
	Userid      string    `gorm:"column:user_id;foreign_key:Userid" json:"userid,omitempty"`
	WatchlistId int      `gorm:"column:watchlist_id" json:"watchlist_id,omitempty"`
	StockName   string    `gorm:"column:stock_name" json:"stock_name,omitempty"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (u *UserWatchlist) TableName() string {
	return "user_watchlist"
}
