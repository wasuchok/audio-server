package handlers

import (
	"net/http"
	"servergo/player"

	"github.com/labstack/echo/v4"
)

func HandleGetMp3ForClient(c echo.Context) error {
	data := player.GetFullBuffer()
	if len(data) == 0 {
		return c.String(http.StatusNotFound, "No audio buffer")
	}

	return c.Blob(http.StatusOK, "audio/mpeg", data)
}
