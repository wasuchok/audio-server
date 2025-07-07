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

func HandleConvertWav(c echo.Context) error {
	var req models.ConvertWavRequest
	if err := c.Bind(&req); err != nil || req.InputPath == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Volume <= 0 {
		req.Volume = 1
	}

	outputName := strings.TrimSuffix(req.OutputName, ".wav")
	if outputName == "" {
		outputName = utils.GenerateUUID()
	}
	outputPath := filepath.Join(outputDir, outputName+".wav")

	cmd := exec.Command("ffmpeg", "-i", req.InputPath, "-acodec", "pcm_s16le", "-ac", "1", "-ar", "44100", "-af", fmt.Sprintf("volume=%.2f", req.Volume), outputPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "ffmpeg failed", "details": string(out)})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "true", "outputWav": outputPath})
}
