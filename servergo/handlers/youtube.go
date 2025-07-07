package handlers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"servergo/models"
	"servergo/utils"

	"github.com/labstack/echo/v4"
)

const outputDir = "output"

func HandleYoutubeToWav(c echo.Context) error {
	var req models.YoutubeRequest
	if err := c.Bind(&req); err != nil || req.YoutubeURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing YouTube URL"})
	}

	filename := strings.TrimSpace(req.Filename)
	if filename == "" {
		filename = utils.GenerateUUID()
	}
	filename = strings.ReplaceAll(filename, " ", "_")
	mp3Path := filepath.Join(outputDir, fmt.Sprintf("%s.mp3", filename))
	wavPath := filepath.Join(outputDir, fmt.Sprintf("%s.wav", filename))

	// 📥 1. ดาวน์โหลด MP3 จาก YouTube
	downloadCmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "-o", mp3Path, req.YoutubeURL)
	out, err := downloadCmd.CombinedOutput()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "yt-dlp failed",
			"details": string(out),
		})
	}

	// 🎛 2. แปลง MP3 → WAV ให้ตรง format เดิม
	convertCmd := exec.Command("ffmpeg",
		"-i", mp3Path,
		"-acodec", "pcm_s16le",
		"-ac", "1",
		"-ar", "44100",
		"-af", "volume=1",
		wavPath,
	)

	out, err = convertCmd.CombinedOutput()

	defer os.Remove(mp3Path)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "ffmpeg failed",
			"details": string(out),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"success": "true",
		"wavFile": wavPath,
	})
}
