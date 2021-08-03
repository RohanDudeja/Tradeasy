package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Watchlist struct {
	gorm.Model
	WatchListId   int    `json:"watchlist_id,omitempty" gorm:"primary_key"`
	WatchListName string `gorm:"watchlist_name" json:"watchlist_name,omitempty"`
}

func (w *Watchlist) TableName() string {
	return "watchlist"
}
