package services

import (
	"Tradeasy/config"
	"Tradeasy/internal/model"
	"Tradeasy/internal/services/order"
	"Tradeasy/internal/services/stock_exchange"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

func OrderConnection() {

	//setting up order connection
	connUrl := "ws://localhost:8080" + "/socket" + "/orders"
	conn, _, err := websocket.DefaultDialer.Dial(connUrl, nil)

	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()
	// Our main loop
	for {
		_, orderMessage, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		var orderDetails stock_exchange.OrderResponse
		err=json.Unmarshal(orderMessage,&orderDetails)
		if err!=nil{
			log.Println("Error during Unmarshalling:",err)
			break
		}
		var p model.PendingOrders
		if err=config.DB.Table("pending_orders").Where("order_id=?",orderDetails.OrderID).First(&p).Error;err!=nil{
			log.Println("Order_id doesnt match with pending_orders table:",err)
			break
		}
		if p.OrderType=="Buy"{
			err=order.UpdateBuyOrder(orderDetails)
			if err!=nil {
				log.Println("Error in updating buy order")
				break
			}
		}else if p.OrderType=="Sell"{
			err=order.UpdateSellOrder(orderDetails)
			if err!=nil {
				log.Println("Error in updating buy order")
				break
			}
		}


		log.Printf("Received: %s", orderMessage)
	}
}
func StockConnection() {

	//setting up stock connection
	socketUrl := "ws://localhost:8080" + "/socket" + "/stocks"
	stockConn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)

	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer stockConn.Close()

	// Our main loop
	for {
		_, stockMessage, err := stockConn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		var stockDetails []stock_exchange.StockDetails
		err=json.Unmarshal(stockMessage,&stockDetails)
		if err!=nil{
			log.Println("Error during Unmarshalling:",err)
			break
		}
		err=order.UpdateStocksFeed(stockDetails)
		if err!=nil {
			log.Println("Error during saving the stocks feed:",err)
			break
		}
		log.Printf("Received: %s", stockMessage)
	}
}
