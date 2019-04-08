package model

import "encoding/xml"

type MediaInfo struct {
	XMLName       xml.Name `xml:"MediaInfo"`
	Medias        []Media  `xml:"media"`
}
