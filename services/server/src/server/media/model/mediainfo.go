package model

import (
	"encoding/xml"
	"time"
	"strings"
	"server/core/repository"
)

type MediaInfo struct {
	repository.Model
	XMLName       xml.Name `xml:"MediaInfo"`
	Medias        []Media  `xml:"media"`
	UserID        uint
	JobID         string
	VideoID       string
	Title         string
	Description   string
	Published     bool
	PublishedDate time.Time
}

func (info *MediaInfo) CanPublish() (bool) {
	return len(strings.TrimSpace(info.Description)) > 0
}
