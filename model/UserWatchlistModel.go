package model

import (
	"github.com/jinzhu/gorm"
	"time"
)
type UserWatchlist struct {
	gorm.Model
	ID                    int       `gorm:"primary_key" json:"id,omitempty"`
	Userid                string   `gorm:"foreign_key:Userid" json:"userid,omitempty"`
	WatchListId           string       `json:"watchlist_id,omitempty"`
	StockName             string       `json:"stock_name,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	DeletedAt             time.Time `json:"deleted_at"`
}

type Watchlist struct {
	gorm.Model
	WatchListId           int    `json:"watchlist_id,omitempty" gorm:"primary_key"`
	WatchListName      string    `json:"watchlist_name,omitempty"`
}



func (u *UserWatchlist) TableName() string{
	return "user_watchlist"
}

func (w *Watchlist) TableName() string{
	return "watchlist"
}

