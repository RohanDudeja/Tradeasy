package services

import (
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
		log.Printf("Received: %s", stockMessage)
	}
}
