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

var ChunkSize = 1024 // เพิ่ม chunk size
var IntervalMs = 5   // ลด interval

func HandleMicWebSocket(w http.ResponseWriter, r *http.Request) {
	player.Pause()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	log.Println("🌐 Mic WebSocket connected")

	// ✨ สร้าง ffmpeg command สำหรับ MP3
	// cmd := exec.Command("ffmpeg",
	// 	"-f", "s16le",
	// 	"-ar", "44100",
	// 	"-ac", "1",
	// 	"-i", "pipe:0",
	// 	"-acodec", "libmp3lame",
	// 	"-ar", "16000", // เปลี่ยนเป็น 16kHz
	// 	"-ac", "1",
	// 	"-b:a", "64k", // ลด bitrate สำหรับ 16kHz
	// 	"-f", "mp3",
	// 	"pipe:1",
	// )

	// ✨ สร้าง ffmpeg command
	// cmd := exec.Command("ffmpeg",
	// 	"-f", "s16le",
	// 	"-ar", "44100",
	// 	"-ac", "1",
	// 	"-i", "pipe:0",
	// 	"-acodec", "libmp3lame",
	// 	"-ar", "16000",
	// 	"-b:a", "64k",
	// 	"-ac", "1",
	// 	"-af", "volume=1",
	// 	"-f", "mp3",
	// 	"pipe:1",
	// )

	cmd := exec.Command("ffmpeg",
		"-f", "s16le", // รูปแบบ input เป็น PCM 16-bit little-endian
		"-ar", "44100", // Sample rate 44.1kHz เพื่อคุณภาพเสียงที่ดี
		"-ac", "1", // โมโน (1 channel)
		"-i", "pipe:0", // อ่าน input จาก stdin
		"-acodec", "libmp3lame", // ใช้ codec MP3
		"-ar", "44100", // Output sample rate 44.1kHz
		"-ac", "2", // ยังคงเป็นโมโน
		"-b:a", "128k", // Bitrate 128kbps เพื่อคุณภาพที่เหมาะสม
		"-af", "volume=1.5", // ปรับ volume 1.5 เท่า ลดการ clipping
		"-f", "mp3", // Output format เป็น MP3
		"pipe:1", // ส่ง output ไปยัง stdout
	)

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
