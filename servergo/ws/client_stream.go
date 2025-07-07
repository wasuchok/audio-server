package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var clientConns = struct {
	sync.RWMutex
	clients map[*websocket.Conn]bool
}{
	clients: make(map[*websocket.Conn]bool),
}

var upgraderClient = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleClientAudioStream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgraderClient.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer func() {
		clientConns.Lock()
		delete(clientConns.clients, conn)
		clientConns.Unlock()
		conn.Close()
	}()

	clientConns.Lock()
	clientConns.clients[conn] = true
	clientConns.Unlock()

	log.Println("🖥️ Web client audio connected")

	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}
	log.Println("🔌 Web client audio disconnected")
}
