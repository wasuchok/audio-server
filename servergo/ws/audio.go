package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	AudioClients   = make(map[*websocket.Conn]bool)
	AudioClientsMu sync.Mutex
)

func HandleAudioWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ Failed to upgrade websocket:", err)
		return
	}

	AudioClientsMu.Lock()
	AudioClients[conn] = true
	AudioClientsMu.Unlock()

	log.Println("🔊 New WebSocket audio client connected")

	go func() {
		defer func() {
			AudioClientsMu.Lock()
			delete(AudioClients, conn)
			AudioClientsMu.Unlock()
			conn.Close()
			log.Println("🔌 WebSocket audio client disconnected")
		}()

		// Keep reading to detect disconnection
		for {
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}()
}
