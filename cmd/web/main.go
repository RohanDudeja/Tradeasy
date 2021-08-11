package main

import (
	"Tradeasy/config"
	"Tradeasy/internal/provider/redis"
	"Tradeasy/internal/router"
	"Tradeasy/internal/services/order"
	"log"
)

func main() {
	//initialise db
	err_ := config.InitialiseDB()
	if err_ != nil {
		log.Fatalf("Gorm: failed to open DB: %v\n", err_)
	}
	redis.CreateClient()
	pong, err := redis.TestClient()
	if err != nil {
		log.Fatalf("Redis: failed to create client: %v\n", err)
	}
	log.Println(pong)

	//setup router
	r := router.SetUpRouter()
	order.InitialiseClientSocket()
	r.Run(config.ServerURL(config.GetConfig()))
}
