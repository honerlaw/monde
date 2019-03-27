package util

import (
	"os/exec"
	"lambda/aws"
	"encoding/xml"
	"errors"
	"server/media/model"
	"log"
	serverAWS "server/core/service/aws"
	"server/core/repository"
)

func GetMediaInfo(metadata *aws.S3RecordMetadata) (*model.MediaInfo, error) {
	url, err := serverAWS.GetS3Service().GetSignedUrl(metadata.Bucket, metadata.Key)

	if err != nil {
		return nil, err
	}

	cmd := exec.Command("bin/mediainfo", "--full", "--output=XML", *url);
	data, err := cmd.Output();
	if err != nil {
		log.Print(string(err.(*exec.ExitError).Stderr))
		return nil, err
	}

	var info model.MediaInfo
	if err = xml.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	info.UserID = metadata.UserId
	info.VideoID = metadata.VideoId

	return &info, nil;
}

func ValidateMediaInfo(mediainfo *model.MediaInfo) (error) {
	for _, mediaInfoMedia := range mediainfo.Medias {
		for _, track := range mediaInfoMedia.Tracks {
			if track.Duration > 30 {
				return errors.New("Duration must be less than 30 seconds")
			}
		}
	}
	return nil
}

func Insert(mediainfo *model.MediaInfo) (error) {
	tx := repository.DB.Begin()

	tx.Save(mediainfo)

	if tx.Error != nil {
		tx.Rollback()
		log.Print(tx.Error)
		return tx.Error
	}

	tx.Commit()

	return nil
}