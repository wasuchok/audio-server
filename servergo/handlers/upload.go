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

func HandleConvertMp3FromFile(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		name = utils.GenerateUUID()
	}
	name = filepath.Base(name)
	rawPath := filepath.Join(outputDir, name+"_raw")      // อัปโหลดมาชื่อ raw
	finalMp3Path := filepath.Join(outputDir, name+".mp3") // หลังแปลงแล้วใช้ชื่อจริง

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing file"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
	}
	defer src.Close()

	dst, err := os.Create(rawPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}

	// ✅ Convert rawPath → finalMp3Path ด้วย ffmpeg
	convertCmd := exec.Command("ffmpeg",
		"-y",
		"-i", rawPath,
		"-acodec", "libmp3lame",
		"-ac", "1",
		"-ar", "44100",
		"-b:a", "16k",
		"-af", "volume=1",
		finalMp3Path,
	)
	out, err := convertCmd.CombinedOutput()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "FFmpeg conversion failed",
			"details": string(out),
		})
	}

	// 🧹 ลบไฟล์ raw หลังแปลงเสร็จ (optional)
	os.Remove(rawPath)

	return c.JSON(http.StatusOK, map[string]string{
		"success":  "true",
		"mp3File":  finalMp3Path,
		"filename": file.Filename,
	})
}
