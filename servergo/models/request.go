package models

type YoutubeRequest struct {
	YoutubeURL string `json:"youtubeUrl"`
	Filename   string `json:"filename"`
}

type ConvertWavRequest struct {
	InputPath  string  `json:"inputPath"`
	OutputName string  `json:"outputName"`
	Volume     float64 `json:"volume"`
}
