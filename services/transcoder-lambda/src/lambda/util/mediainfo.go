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
	"fmt"
)

func GetMediaInfo(metadata *aws.S3RecordMetadata) (*model.MediaInfo, error) {
	info := model.MediaInfo{}
	info.UserID = metadata.UserId
	info.VideoID = metadata.VideoId

	url, err := serverAWS.GetS3Service().GetSignedUrl(metadata.Bucket, metadata.Key)

	if err != nil {
		return &info, err
	}

	cmd := exec.Command("bin/mediainfo", "--full", "--output=XML", *url);
	data, err := cmd.Output();
	if err != nil {
		log.Print(string(err.(*exec.ExitError).Stderr))
		return &info, err
	}

	if err = xml.Unmarshal(data, &info); err != nil {
		return &info, err
	}

	return &info, nil;
}

func ValidateMediaInfo(mediainfo *model.MediaInfo) (error) {
	if len(mediainfo.Medias) == 0 {
		return errors.New("no media found for file")
	}

	for _, mediaInfoMedia := range mediainfo.Medias {

		if len(mediaInfoMedia.Tracks) == 0 {
			return errors.New("no tracks found for file")
		}

		videoTrackCount := 0
		for _, track := range mediaInfoMedia.Tracks {
			// make sure it has a video track
			if track.Type == "Video" {
				videoTrackCount += 1

				// and that video track is 300 seconds or less
				if track.Duration > 300 {
					return errors.New("duration must be less than 5 minutes")
				}
			}
		}

		// we should only have one video track per file
		// @todo figure out cases where there is multiple video tracks...
		if videoTrackCount != 1 {
			return errors.New(fmt.Sprintf("incorrect number of video tracks found: %d", videoTrackCount))
		}
	}
	return nil
}

func RetryInsert(mediainfo *model.MediaInfo, retries uint) {
	log.Printf("attempting to insert media info, retries left: %d", retries)
	if retries <= 0 {
		return
	}

	err := Insert(mediainfo)

	if err != nil {
		log.Print("Failed to insert media and job information", err)
		RetryInsert(mediainfo, retries - 1)
	}
}

func Insert(mediainfo *model.MediaInfo) (error) {
	// check if there is a db connection, if not create one
	err := repository.DB.DB().Ping()

	if err != nil {
		log.Print(err);
		return err;
	}

	tx := repository.DB.Begin()

	// @todo this is generating an update for the tracks, when we should always insert...
	tx.Create(mediainfo)

	if tx.Error != nil {
		tx.Rollback()
		log.Print(tx.Error)
		return tx.Error
	}

	tx.Commit()

	return nil
}
