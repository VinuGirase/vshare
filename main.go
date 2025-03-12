package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}

		for client := range clients {
			if client != conn {
				client.WriteMessage(websocket.BinaryMessage, msg)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws/1", handleConnections)

	port := ":8080"
	fmt.Println("Server running on http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
