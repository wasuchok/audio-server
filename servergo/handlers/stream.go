package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

var browserClients = make(map[chan []byte]bool)

func StreamMP3ToBrowser(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "audio/mpeg")
	c.Response().WriteHeader(http.StatusOK)

	writer := c.Response().Writer
	flusher, ok := writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming not supported")
	}

	clientChan := make(chan []byte, 100)
	browserClients[clientChan] = true
	defer func() {
		delete(browserClients, clientChan)
		close(clientChan)
	}()

	notify := c.Request().Context().Done()

	for {
		select {
		case chunk := <-clientChan:
			_, err := writer.Write(chunk)
			if err != nil {
				log.Println("❌ Failed writing to browser:", err)
				return nil
			}
			flusher.Flush()
		case <-notify:
			log.Println("🛑 Browser disconnected from /live")
			return nil
		}
	}
}

func BroadcastToBrowsers(chunk []byte) {
	for ch := range browserClients {
		select {
		case ch <- chunk:
		default:
			log.Println("⚠️ Browser client buffer full, skipping chunk")
		}
	}
}
