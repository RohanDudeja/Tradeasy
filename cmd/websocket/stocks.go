package main

import (
	"github.com/gorilla/websocket"
	"log"
)

func main() {

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
