package server

import (
	"fmt"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

var ESPConns = make(map[string]net.Conn) // key = IP หรือ MAC
var WebClients = make(map[*websocket.Conn]bool)

func StartTCPServer(port int) {
	address := fmt.Sprintf("0.0.0.0:%d", port)
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

			tcpConn, ok := conn.(*net.TCPConn)
			if ok {
				tcpConn.SetKeepAlive(true)
				tcpConn.SetKeepAlivePeriod(30 * time.Second)
			}

			ip := conn.RemoteAddr().String()
			fmt.Printf("✅ ESP32 connected from %s\n", ip)

			// ถ้ามีการเชื่อมต่อซ้ำเดิม ให้ปิดตัวเก่า
			if oldConn, exists := ESPConns[ip]; exists {
				fmt.Printf("♻️  Closing old connection from %s\n", ip)
				oldConn.Close()
			}

			ESPConns[ip] = conn
			fmt.Println("🎵 Audio output ready - ESP32 can now receive audio data")

			go handleTCPConnection(conn, ip)
		}
	}()
}

func handleTCPConnection(conn net.Conn, ip string) {
	defer func() {
		fmt.Printf("🔌 ESP32 disconnected from %s\n", ip)
		conn.Close()
		delete(ESPConns, ip)
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("⚠️ Error reading from %s: %v\n", ip, err)
			break
		}

		_ = buf[:n]
		// Optional: log or process buf[:n]
	}
}
