package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"servergo/handlers"
	"servergo/player"
	"servergo/server"
	"servergo/ws"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const outputDir = "output"

var playlist = []string{
	"output/song2.mp3",
	"output/song1.mp3",
	"output/song5.mp3",
	"output/song6.mp3",
	"output/song3.mp3",
	"output/song4.mp3",
}

var currentTrackIndex = 0

var (
// ลบตัวแปร SendWavHeader ออก
)

func loadAndPlayCurrentTrack() {
	// 🔁 วนกลับไปเพลงแรกถ้าจบ playlist
	if currentTrackIndex >= len(playlist) {
		currentTrackIndex = 0
	}

	file := playlist[currentTrackIndex]
	log.Printf("📂 Loading file: %s", file)

	data, err := os.ReadFile(file)
	if err != nil {
		log.Printf("❌ Failed to load %s: %v", file, err)
		return
	}

	log.Printf("📦 File loaded: %s (%d bytes)", file, len(data))

	if len(data) == 0 {
		log.Printf("❌ File is empty: %s", file)
		return
	}

	player.SetBuffer(data)

	if server.ESPConn == nil {
		log.Println("⚠️ ESP32 not connected - audio will not play")
	} else {
		log.Println("✅ ESP32 connected - ready to play audio")
	}

	player.Play()
	log.Printf("📦 Loaded %s (%d bytes)", file, len(data))

	log.Println("🎶 Now playing:", file)
}

func main() {
	os.MkdirAll(outputDir, os.ModePerm)

	http.HandleFunc("/ws/mic", ws.HandleMicWebSocket)
	http.HandleFunc("/ws/control", ws.HandleControlWebSocket)
	http.HandleFunc("/ws/set-mic-volume", handlers.HandleSetMicVolume)
	http.HandleFunc("/ws/audio", ws.HandleAudioWebSocket)

	go func() {
		log.Println("🌐 WebSocket listening on :7777 at /ws/mic")
		err := http.ListenAndServe(":7777", nil)
		if err != nil {
			log.Fatal("WebSocket server error:", err)
		}
	}()

	player.SendChunk = func(chunk []byte) {
		// ส่งให้ ESP32
		if server.ESPConn != nil {
			server.ESPConn.Write(chunk)
		}

		// 🔊 ส่งให้เบราว์เซอร์
		handlers.BroadcastToBrowsers(chunk)

		time.Sleep(10 * time.Millisecond)
	}

	// player.SendChunk = func(chunk []byte) {
	// 	// 1. ส่งให้ ESP32
	// 	if server.ESPConn != nil {
	// 		_, err := server.ESPConn.Write(chunk)
	// 		if err != nil {
	// 			log.Println("❌ ESP32 error:", err)
	// 			server.ESPConn = nil
	// 		}
	// 	}

	// 	// 2. ส่งให้ WebSocket Clients
	// 	ws.AudioClientsMu.Lock()
	// 	for client := range ws.AudioClients {
	// 		err := client.WriteMessage(websocket.BinaryMessage, chunk)
	// 		if err != nil {
	// 			log.Println("❌ Error sending to WebSocket client:", err)
	// 			client.Close()
	// 			delete(ws.AudioClients, client)
	// 		}
	// 	}
	// 	ws.AudioClientsMu.Unlock()

	// 	time.Sleep(10 * time.Millisecond)
	// }

	// player.SendChunk = func(chunk []byte) {
	// 	if server.ESPConn != nil {
	// 		log.Printf("📤 Sending %d bytes to ESP32", len(chunk))

	// 		_, err := server.ESPConn.Write(chunk)
	// 		if err != nil {
	// 			log.Println("❌ Failed to send chunk to ESP32:", err)
	// 			server.ESPConn = nil
	// 			return
	// 		}

	// 		if flusher, ok := server.ESPConn.(interface{ Flush() error }); ok {
	// 			flusher.Flush()
	// 		}

	// 		log.Printf("✅ Successfully sent %d bytes to ESP32", len(chunk))

	// 		time.Sleep(10 * time.Millisecond)
	// 	} else {
	// 		log.Println("⚠️ ESP32 not connected - audio chunks not being sent")
	// 	}
	// }

	loadAndPlayCurrentTrack()

	player.OnFinishTrack = func() {
		currentTrackIndex++
		if currentTrackIndex >= len(playlist) {
			currentTrackIndex = 0
		}
		loadAndPlayCurrentTrack()
	}

	go server.StartTCPServer(5555)

	// แสดงสถานะการเชื่อมต่อ
	go func() {
		for {
			time.Sleep(5 * time.Second)
			if server.ESPConn == nil {
				log.Println("⚠️ ESP32 not connected - waiting for connection on port 5555")
			} else {
				log.Println("✅ ESP32 connected - audio ready")
			}
		}
	}()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	}))
	e.POST("/api/youtube-to-mp3", handlers.HandleYoutubeToMp3)
	e.POST("/api/convert-mp3", handlers.HandleConvertMp3)
	e.POST("/api/upload-mp3", handlers.HandleConvertMp3FromFile)
	e.GET("/api/audio.mp3", handlers.HandleGetMp3ForClient)
	e.GET("/live", handlers.StreamMP3ToBrowser)

	fmt.Println("🚀 API running on http://localhost:4444")
	e.Logger.Fatal(e.Start(":4444"))
}
