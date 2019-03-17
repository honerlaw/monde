package util

import (
	"os/exec"
	"lambda/aws"
	"encoding/xml"
	"errors"
	"package/model/media"
	"log"
)

func GetMediaInfo(metadata *aws.S3RecordMetadata) (*media.MediaInfo, error) {
	url, err := aws.GetSignedS3Url(metadata.Bucket, metadata.Key)

	if err != nil {
		return nil, err
	}

	cmd := exec.Command("bin/mediainfo", "--full", "--output=XML", *url);
	data, err := cmd.Output();
	if err != nil {
		log.Print(string(err.(*exec.ExitError).Stderr))
		return nil, err
	}

	var info media.MediaInfo
	if err = xml.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	info.UserID = metadata.UserId
	info.VideoID = metadata.VideoId

	return &info, nil;
}

func ValidateMediaInfo(mediainfo *media.MediaInfo) (error) {
	for _, media := range mediainfo.Medias {
		for _, track := range media.Tracks {
			if track.Duration > 30 {
				return errors.New("Duration must be less than 30 seconds")
			}
		}
	}
	return nil
}
