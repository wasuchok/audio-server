package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"servergo/handlers"
	"servergo/player"
	"servergo/server"
	"servergo/ws"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const outputDir = "output"

var playlist = []string{
	"output/Always.wav",
	"output/intro.wav",
	"output/Sanctuary.wav",
	"output/Stuck_with_U.wav",
	"output/deja_vu.wav",
	"output/party_4_u.wav",
}

var currentTrackIndex = 0

func loadAndPlayCurrentTrack() {
	// 🔁 วนกลับไปเพลงแรกถ้าจบ playlist
	if currentTrackIndex >= len(playlist) {
		currentTrackIndex = 0
	}

	file := playlist[currentTrackIndex]
	data, err := os.ReadFile(file)
	if err != nil {
		log.Printf("❌ Failed to load %s: %v", file, err)
		return
	}

	player.SetBuffer(data)
	player.Play()
	log.Printf("📦 Loaded %s (%d bytes)", file, len(data))

	log.Println("🎶 Now playing:", file)
}

func main() {
	os.MkdirAll(outputDir, os.ModePerm)

	http.HandleFunc("/ws/mic", ws.HandleMicWebSocket)
	http.HandleFunc("/ws/control", ws.HandleControlWebSocket)
	http.HandleFunc("/ws/stream-client", ws.HandleMicWebSocketForWeb)

	go func() {
		log.Println("🌐 WebSocket listening on :7777 at /ws/mic")
		err := http.ListenAndServe(":7777", nil)
		if err != nil {
			log.Fatal("WebSocket server error:", err)
		}
	}()

	player.SendChunk = func(chunk []byte) {
		if server.ESPConn != nil {
			_, err := server.ESPConn.Write(chunk)
			if err != nil {
				log.Println("❌ Failed to send chunk to ESP32:", err)
			}

		}

		ws.BroadcastToClients(player.MakeWavChunk(chunk))
	}

	loadAndPlayCurrentTrack()

	player.OnFinishTrack = func() {
		currentTrackIndex++
		if currentTrackIndex >= len(playlist) {
			currentTrackIndex = 0
		}
		loadAndPlayCurrentTrack()
	}

	go server.StartTCPServer(5555)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	}))
	e.POST("/api/youtube-to-wav", handlers.HandleYoutubeToWav)
	e.POST("/api/convert-wav", handlers.HandleConvertWav)
	e.POST("/api/upload-wav", handlers.HandleConvertWavFromFile)
	e.GET("/api/audio.wav", handlers.HandleGetWavForClient)

	fmt.Println("🚀 API running on http://localhost:4444")
	e.Logger.Fatal(e.Start(":4444"))
}
