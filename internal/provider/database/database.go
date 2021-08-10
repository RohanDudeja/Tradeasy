package database

import (
	"Tradeasy/config"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// InitialiseDB ...assign connection to global *gorm.DB variable DB
func InitialiseDB() error {
	dbString := config.DbURL(config.GetConfig())
	var err error
	db, err = gorm.Open("mysql", dbString)
	if err != nil {
		return err
	}
	if gin.IsDebugging() {
		db.LogMode(true)
	}
	return nil
}
func GetDB() *gorm.DB {
	return db
}

func SetDB(database *gorm.DB) {
	db = database
}
