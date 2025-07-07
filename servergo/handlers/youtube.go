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

func HandleYoutubeToMp3(c echo.Context) error {
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

	// 📥 ดาวน์โหลด MP3 จาก YouTube
	downloadCmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0", // คุณภาพสูงสุด
		"-o", mp3Path,
		req.YoutubeURL)
	out, err := downloadCmd.CombinedOutput()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "yt-dlp failed",
			"details": string(out),
		})
	}

	// 🎛 แปลง MP3 ให้ตรง format ที่ต้องการ
	convertCmd := exec.Command("ffmpeg",
		"-i", mp3Path,
		"-acodec", "libmp3lame",
		"-ac", "1",
		"-ar", "44100",
		"-b:a", "16k",
		"-af", "volume=1",
		mp3Path+".tmp",
	)

	out, err = convertCmd.CombinedOutput()

	if err != nil {
		os.Remove(mp3Path)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "ffmpeg failed",
			"details": string(out),
		})
	}

	// ย้ายไฟล์ชั่วคราวมาแทนที่ไฟล์เดิม
	os.Remove(mp3Path)
	os.Rename(mp3Path+".tmp", mp3Path)

	return c.JSON(http.StatusOK, map[string]string{
		"success": "true",
		"mp3File": mp3Path,
	})
}
