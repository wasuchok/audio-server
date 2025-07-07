package server

import (
	"fmt"
	"net"

	"github.com/gorilla/websocket"
)

var ESPConn net.Conn
var WebClients = make(map[*websocket.Conn]bool)

func StartTCPServer(port int) {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to start TCP server: %v", err))
	}
	fmt.Printf("📡 TCP Server listening on port %d\n", port)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("❌ Failed to accept TCP connection:", err)
				continue
			}

			fmt.Println("✅ ESP32 connected via TCP")
			ESPConn = conn
			fmt.Println("🎵 Audio output ready - ESP32 can now receive audio data")

			go handleTCPConnection(conn)
		}
	}()
}

func handleTCPConnection(conn net.Conn) {
	defer func() {
		fmt.Println("🔌 ESP32 disconnected")
		ESPConn = nil
		conn.Close()
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("⚠️ Error reading from ESP32:", err)
			break
		}
		// Optional: log incoming data
		_ = buf[:n]
	}
}
