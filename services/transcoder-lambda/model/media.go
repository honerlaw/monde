package model

import "encoding/xml"

type Media struct {
	XMLName     xml.Name `xml:"media"`
	Tracks      []Track  `xml:"track"`
}
