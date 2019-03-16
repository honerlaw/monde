package util

import (
	"os/exec"
	"lambda/aws"
	"encoding/xml"
	"fmt"
	"errors"
)

type MediaInfo struct {
	XMLName xml.Name `xml:"MediaInfo"`
	Medias  []Media  `xml:"media"`
}

type Media struct {
	XMLName xml.Name `xml:"media"`
	Tracks  []Track  `xml:"track"`
}

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

func GetMediaInfo(bucket string, key string) (*MediaInfo, error) {
	url, err := aws.GetSignedS3Url(bucket, key)

	if err != nil {
		return nil, err
	}

	cmd := exec.Command("bin/mediainfo", "--full", "--output=XML", *url);
	data, err := cmd.Output();
	if err != nil {
		fmt.Print(string(err.(*exec.ExitError).Stderr))
		return nil, err
	}

	var info MediaInfo
	if err = xml.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	return &info, nil;
}

func ValidateMediaInfo(mediainfo *MediaInfo) (error) {
	for _, media := range mediainfo.Medias {
		for _, track := range media.Tracks {
			if track.Duration > 30 {
				return errors.New("Duration must be less than 30 seconds")
			}
		}
	}
	return nil
}
