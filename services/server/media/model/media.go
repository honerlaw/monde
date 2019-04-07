package model

import (
	"encoding/xml"
	"services/server/core/repository"
)

type Media struct {
	repository.Model
	XMLName     xml.Name `xml:"media"`
	Tracks      []Track  `xml:"track"`
	MediaInfoID uint
}
