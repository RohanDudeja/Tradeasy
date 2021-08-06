package main

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
)

var done chan interface{}
var interrupt chan os.Signal

func main() {

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
