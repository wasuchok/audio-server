package ws

import (
	"log"
	"net/http"
	"servergo/player"

	"github.com/gorilla/websocket"
)

var upgraderControl = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleControlWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgraderControl.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	log.Println("🎧 Control WebSocket connected")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Control WebSocket error:", err)
			break
		}

		command := string(message)
		log.Println("📥 Command received:", command)

		switch command {
		case "play":
			player.Play()
		case "pause":
			player.Pause()
		case "resume":
			player.Resume()
		case "stop":
			player.Stop()
		default:
			log.Println("⚠️ Unknown command:", command)
		}
	}

	log.Println("❌ Control WebSocket disconnected")
}
