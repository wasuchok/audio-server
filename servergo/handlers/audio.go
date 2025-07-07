package handlers

import (
	"net/http"
	"servergo/player"
	"strconv"

	"github.com/labstack/echo/v4"
)

var micVolume = 1.0 // ค่าเริ่มต้น

func HandleGetMp3ForClient(c echo.Context) error {
	data := player.GetFullBuffer()
	if len(data) == 0 {
		return c.String(http.StatusNotFound, "No audio buffer")
	}

	return c.Blob(http.StatusOK, "audio/mpeg", data)
}

func HandleSetMicVolume(w http.ResponseWriter, r *http.Request) {
	vStr := r.URL.Query().Get("v")
	v, err := strconv.ParseFloat(vStr, 64)
	if err != nil {
		http.Error(w, "invalid volume", 400)
		return
	}
	SetMicVolume(v)
	w.Write([]byte("ok"))
}

func SetMicVolume(v float64) {
	micVolume = v
}
