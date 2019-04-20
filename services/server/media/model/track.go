package model

import (
	"services/server/core/repository"
)

type Track struct {
	repository.Model
	MediaID      string  `json:"media_id" column:"media_id"`
	Type         string  `json:"type" column:"type"`
	Duration     float64 `json:"duration" column:"duration"`
	Width        int64   `json:"width" column:"width"`
	Height       int64   `json:"height" column:"height"`
	Format       string  `json:"format" column:"format"`
	Encoded_Date string  `json:"encoded_date" column:"encoded_date"`
	VideoCount   string  `json:"video_count" column:"video_count"`
	DataSize     int64   `json:"data_size" column:"data_size"`
	FileSize     int64   `json:"file_size" column:"file_size"`
}
