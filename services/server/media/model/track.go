package model

import (
	"services/server/core/repository"
)

type Track struct {
	repository.Model
	MediaID      string
	Type         string
	Duration     float64
	Width        int64
	Height       int64
	Format       string
	Encoded_Date string
	VideoCount   string
	DataSize     int64
	FileSize     int64
}
