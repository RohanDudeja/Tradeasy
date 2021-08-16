package order

import (
	"Tradeasy/config"
	_ "Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/provider/database"
	"Tradeasy/internal/services/stock_exchange"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func OrderConnection() {

	//setting up order connection
	connUrl := "ws://localhost:8080" + "/socket" + "/orders"

	userName := config.GetConfig().StockExchange.Authentication.UserName
	password := config.GetConfig().StockExchange.Authentication.Password

	h := http.Header{"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(userName+":"+password))}}

	conn, _, err := websocket.DefaultDialer.Dial(connUrl, h)

	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()
	// Our main loop
	for {
		_, orderMessage, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			continue
		}
		var orderDetails stock_exchange.OrderResponse
		err = json.Unmarshal(orderMessage, &orderDetails)
		if err != nil {
			log.Println("Error during Unmarshalling:", err)
			continue
		}
		var p model.PendingOrders
		if err = database.GetDB().Table("pending_orders").Where("order_id=?", orderDetails.OrderID).First(&p).Error; err != nil {
			log.Println("Order_id doesnt match with pending_orders table:", err)
			continue
		}
		if p.OrderType == "Buy" {
			go func() {
				UpdateBuyOrder(&orderDetails)
			}()
		} else if p.OrderType == "Sell" {
			go func() {
				UpdateSellOrder(&orderDetails)
			}()
		}
		log.Printf("Received the order details from Stock Exchange Engine: %s", orderMessage)
	}
}
func StockConnection() {

	//setting up stock connection
	socketUrl := "ws://localhost:8080" + "/socket" + "/stocks"
	userName := config.GetConfig().StockExchange.Authentication.UserName
	password := config.GetConfig().StockExchange.Authentication.Password

	h := http.Header{"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(userName+":"+password))}}

	stockConn, _, err := websocket.DefaultDialer.Dial(socketUrl, h)

	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer stockConn.Close()

	// Our main loop
	for {
		_, stockMessage, err := stockConn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			continue
		}
		var stockDetails []stock_exchange.StockDetails
		err = json.Unmarshal(stockMessage, &stockDetails)
		if err != nil {
			log.Println("Error during Unmarshalling:", err)
			continue
		}
		go func() {
			UpdateStocksFeed(stockDetails)
		}()
		log.Printf("Received the stock details from Stock Exchange Engine: %s", stockMessage)
	}
}

func InitialiseClientSocket() {
	go OrderConnection()
	go StockConnection()
}
