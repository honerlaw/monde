package model

import (
	"encoding/xml"
)

type Track struct {
	XMLName      xml.Name `xml:"track"`
	Type         string   `xml:"type,attr"`
	Duration     float64  `xml:"Duration"`
	Width        int64    `xml:"Width"`
	Height       int64    `xml:"Height"`
	Format       string   `xml:"Format"`
	Encoded_Date string   `xml:"Encoded_Date"`
	VideoCount   string   `xml:"VideoCount"`
	DataSize     int64    `xml:"DataSize"`
	FileSize     int64    `xml:"FileSize"`
}
