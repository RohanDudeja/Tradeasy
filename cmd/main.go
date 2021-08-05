package cmd

import (
	"Tradeasy/config"
	"Tradeasy/internal/router"
	"github.com/jinzhu/gorm"
	"log"
)

var err error

// initialiseDB ...assign connection to global *gorm.DB variable DB
func initialiseDB() error{
	dbString := config.DbURL(config.BuildConfig())
	config.DB,err = gorm.Open("mysql", dbString)
	if err != nil {
		return err
	}
	return nil
}

func main()  {
	//initialise db
	err_ := initialiseDB()
	if err_ != nil {
		log.Fatalf("Gorm: failed to open DB: %v\n", err_)
	}
	//setup router
	r:=router.SetUpRouter()
	r.Run()
}
