package cmd

import (
	"Tradeasy/config"
	"Tradeasy/internal/router"
	"github.com/jinzhu/gorm"
	"log"
)

var err error

func main()  {

	//assign connection to global *gorm.DB variable DB
	dbstring := config.DbURL(config.BuildConfig())
	config.DB,err = gorm.Open("mysql", dbstring)
	if err != nil {
		log.Fatalf("Gorm: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := config.DB.Close(); err != nil {
			log.Fatalf("Gorm: failed to close DB: %v\n", err)
		}
	}()

	r:=router.SetUpRouter()
	r.Run()
}
