package main

import (
	"Tradeasy/config"
	"Tradeasy/internal/router"
	"log"
)

func main() {
	//initialise db
	err_ := config.InitialiseDB()
	if err_ != nil {
		log.Fatalf("Gorm: failed to open DB: %v\n", err_)
	}
	//setup router
	r := router.SetUpRouter()
	r.Run(config.ServerURL(config.Con))
}
