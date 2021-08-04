package model

import (
	"time"
)

type Watchlist struct {
	Id        int       `json:"watchlist_id,omitempty" gorm:"column:id;primary_key"`
	Name      string    `gorm:"column:name" json:"watchlist_name,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

func (w *Watchlist) TableName() string {
	return "watchlist"
}
