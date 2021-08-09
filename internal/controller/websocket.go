package controller

import (
	"Tradeasy/internal/services/stock_exchange"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const TimeInterval = time.Duration(5)*time.Millisecond*1000

func StockHandler(c *gin.Context) {
	var upgrader = websocket.Upgrader{} // use default options
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()
	//write stock details
	for range time.Tick(TimeInterval) {
		stocks, err := stock_exchange.StockWrite()
		stockJson, err := json.Marshal(&stocks)
		if err != nil {
			log.Println("Error while converting stocks to bytes", err)
			return
		}

		if err := conn.WriteMessage(websocket.TextMessage, stockJson); err != nil {
			log.Println("Error during writing stocks to websocket:", err)
			return
		}
	}
}
func OrderHandler(c *gin.Context) {
	var upgrader = websocket.Upgrader{} // use default options
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	for {
		select {
		case orderMsg := <-stock_exchange.OrderUpdated:
			orderJson, err := json.Marshal(orderMsg)
			if err != nil {
				log.Println("Error while converting stocks to bytes", err)
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, orderJson); err != nil {
				log.Println("Error during writing stocks to websocket:", err)
				return
			}
		}
	}
}