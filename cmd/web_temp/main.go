package main

import (
	"Tradeasy/config"
	"Tradeasy/internal/router"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	//initialise db
	err_ := config.InitialiseDB()
	if err_ != nil {
		log.Fatalf("Gorm: failed to open DB: %v\n", err_)
	}
	//setup router

	//stock_exchange.InitialiseAllStocks()
	//stock_exchange.InitialiseBuyersAndSellers()
	//stock_exchange.RandomizerAlgo()
	r := router.SetUpRouter()
	if gin.IsDebugging() {
		config.DB.LogMode(true)
	}
	r.Run("localhost:8080")
}
