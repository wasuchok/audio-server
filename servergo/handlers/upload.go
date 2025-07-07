package handlers

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"servergo/utils"

	"github.com/labstack/echo/v4"
)

func HandleConvertWavFromFile(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		name = utils.GenerateUUID()
	}
	name = filepath.Base(name)
	rawWavPath := filepath.Join(outputDir, name+"_raw.wav") // อัปโหลดมาชื่อ raw
	finalWavPath := filepath.Join(outputDir, name+".wav")   // หลังแปลงแล้วใช้ชื่อจริง

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing file"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
	}
	defer src.Close()

	dst, err := os.Create(rawWavPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}

	// ✅ Convert rawWavPath → finalWavPath ด้วย ffmpeg
	convertCmd := exec.Command("ffmpeg",
		"-y",
		"-i", rawWavPath,
		"-acodec", "pcm_s16le",
		"-ac", "1",
		"-ar", "44100",
		"-af", "volume=1",
		finalWavPath,
	)
	out, err := convertCmd.CombinedOutput()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "FFmpeg conversion failed",
			"details": string(out),
		})
	}

	// 🧹 ลบไฟล์ raw หลังแปลงเสร็จ (optional)
	os.Remove(rawWavPath)

	return c.JSON(http.StatusOK, map[string]string{
		"success":  "true",
		"wavFile":  finalWavPath,
		"filename": file.Filename,
	})
}
