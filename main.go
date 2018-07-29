package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"github.com/satori/go.uuid"
	"pubsub/pubsub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func autoId() (string) {

	return uuid.Must(uuid.NewV4()).String()
}

var ps = &pubsub.PubSub{}

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true

	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}


	client := pubsub.Client{
		Id:         autoId(),
		Connection: conn,
	}

	// add this client into the list
	ps.AddClient(client)

	fmt.Println("New Client is connected, total: ", len(ps.Clients))

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Something went wrong", err)

			ps.RemoveClient(client)
			log.Println("total clients and subscriptions ", len(ps.Clients), len(ps.Subscriptions))


			return
		}

		ps.HandleReceiveMessage(client, messageType, p)

	}

}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "static")

	})

	http.HandleFunc("/ws", websocketHandler)

	http.ListenAndServe(":3000", nil)

	fmt.Println("Server is running: http://localhost:3000")

}
