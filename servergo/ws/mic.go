package ws

import (
	"io"
	"log"
	"net/http"
	"os/exec"
	"servergo/player"
	"servergo/server"

	"github.com/gorilla/websocket"
)

// 🔁 เพิ่มไว้เพื่อส่งให้ทุก client

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Clients = make(map[*websocket.Conn]bool)

func HandleMicWebSocket(w http.ResponseWriter, r *http.Request) {
	player.Pause()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	log.Println("🌐 Mic WebSocket connected")

	// ✨ สร้าง ffmpeg command
	cmd := exec.Command("ffmpeg",
		"-f", "s16le",
		"-ar", "44100",
		"-ac", "1",
		"-i", "pipe:0",
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "1",
		"-af", "volume=1",
		"-f", "wav",
		"pipe:1",
	)

	// cmd := exec.Command("ffmpeg",
	// 	"-f", "s16le",
	// 	"-ar", "44100",
	// 	"-ac", "1",
	// 	"-i", "pipe:0",

	// 	"-acodec", "pcm_s16le",
	// 	"-ar", "44100",
	// 	"-ac", "1",
	// 	"-af", "volume=1",
	// 	"-f", "wav",
	// 	"-fflags", "nobuffer",
	// 	"-flags", "low_delay",
	// 	"-flush_packets", "1",
	// 	"pipe:1",
	// )

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println("Failed to create stdin pipe:", err)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Failed to create stdout pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println("Failed to start ffmpeg:", err)
		return
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println("ffmpeg stdout read error:", err)
				}
				break
			}

			if server.ESPConn != nil {
				_, err := server.ESPConn.Write(buf[:n])
				if err != nil {
					log.Println("❌ Failed to send to ESP32:", err)
				}
			}
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ WebSocket read error:", err)
			break
		}

		_, err = stdin.Write(data)
		if err != nil {
			log.Println("❌ Failed to write to ffmpeg stdin:", err)
			break
		}
	}

	stdin.Close()
	cmd.Wait()

	log.Println("❌ Mic WebSocket disconnected")
}

func HandleMicWebSocketForWeb(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	Clients[conn] = true
	defer delete(Clients, conn)

	log.Println("🎤 Web mic connected")

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("Mic WS error:", err)
			break
		}

		// 📡 ส่งไปให้ client stream ทุกตัว
		for client := range Clients {
			if err := client.WriteMessage(websocket.BinaryMessage, data); err != nil {
				log.Println("Send to client failed:", err)
			}
		}
	}
}
