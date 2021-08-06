package controller

import (
	"Tradeasy/internal/services/stock_exchange"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var upgrader = websocket.Upgrader{} // use default options
var OrderUpdated = make(chan stock_exchange.OrderResponse)

func StockHandler(c *gin.Context) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()
	//write stock details
	writeEvery(time.Duration(5)*time.Millisecond*1000, conn, stock_exchange.StockWrite)
}
func OrderHandler(c *gin.Context) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	/*go func() {
		for {
			time.Sleep(5 * time.Second)
			OrderUpdated <- stock_exchange.OrderResponse{AveragePrice: uint(rand.Intn(10))}
		}
	}()*/

	for {
		select {
		case orderMsg := <-OrderUpdated:
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
func Home(c *gin.Context) {
	fmt.Fprintf(c.Writer, "Index Page")
}

func writeEvery(d time.Duration, conn *websocket.Conn, f func() (*stock_exchange.StockDetails, error)) {
	for range time.Tick(d) {
		stocks, err := f()
		stockJson, err := json.Marshal(stocks)
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