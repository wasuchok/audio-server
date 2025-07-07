package handlers

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"servergo/models"
	"servergo/utils"

	"github.com/labstack/echo/v4"
)

func HandleConvertMp3(c echo.Context) error {
	var req models.ConvertWavRequest
	if err := c.Bind(&req); err != nil || req.InputPath == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Volume <= 0 {
		req.Volume = 1
	}

	outputName := strings.TrimSuffix(req.OutputName, ".mp3")
	if outputName == "" {
		outputName = utils.GenerateUUID()
	}
	outputPath := filepath.Join(outputDir, outputName+".mp3")

	cmd := exec.Command("ffmpeg",
		"-i", req.InputPath,
		"-acodec", "libmp3lame",
		"-ac", "1",
		"-ar", "44100",
		"-b:a", "16k",
		"-af", fmt.Sprintf("volume=%.2f", req.Volume),
		outputPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ffmpeg failed", "details": string(out)})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "true", "outputMp3": outputPath})
}
