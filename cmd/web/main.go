package main

import (
	"Tradeasy/config"
	"Tradeasy/internal/provider/database"
	"Tradeasy/internal/provider/redis"
	"Tradeasy/internal/router"
	"Tradeasy/internal/services/order"
	"Tradeasy/internal/services/stock_exchange"
	"log"
)

func main() {
	//initialise config
	config.BuildConfig()
	//initialise db
	err_ := database.InitialiseDB()
	if err_ != nil {
		log.Fatalf("Gorm: failed to open DB: %v\n", err_)
	}
	//Initializing Redis Client
	redis.CreateClient()

	//Initializing stocks
	stock_exchange.InitialiseAllStocks()

	//Initializing Randomizer
	stock_exchange.InitialiseRandomizer()

	//setup router
	r := router.SetUpRouter()
	order.InitialiseClientSocket()
	r.Run(config.ServerURL(config.GetConfig()))
}
