package controller

import (
	"Tradeasy/internal/services/stock_exchange"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const stockTimeInterval = time.Duration(5) * time.Second

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
	stock_exchange.GetStockUpdates(conn, stockTimeInterval)
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
	stock_exchange.GetOrderUpdates(conn)
}
