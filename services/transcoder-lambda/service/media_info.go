package service

import (
	"os/exec"
	"services/transcoder-lambda/aws"
	"encoding/xml"
	"errors"
	"log"
	serverAWS "services/server/core/service/aws"
	"services/server/core/repository"
	"fmt"
	"services/transcoder-lambda/model"
	mediaModel "services/server/media/model"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
)

func GetMediaInfo(metadata *aws.S3RecordMetadata) (*model.MediaInfo, error) {
	info := model.MediaInfo{}

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

func RetryInsert(metadata *aws.S3RecordMetadata, mediainfo *model.MediaInfo, job *elastictranscoder.Job, retries uint) {
	log.Printf("attempting to insert media info, retries left: %d", retries)
	if retries <= 0 {
		return
	}

	err := Insert(metadata, mediainfo, job)

	if err != nil {
		log.Print("Failed to insert media and job information", err)
		RetryInsert(metadata, mediainfo, job, retries-1)
	}
}

func Insert(metadata *aws.S3RecordMetadata, mediainfo *model.MediaInfo, job *elastictranscoder.Job) (error) {
	// check if there is a db connection, if not create one
	err := repository.GetRepository().DB().Ping()

	if err != nil {
		log.Print(err);
		return err;
	}

	tx, err := repository.GetRepository().DB().Begin()
	if err != nil {
		log.Fatal(err)
	}

	media := &mediaModel.Media{
		ChannelID: metadata.ChannelID,
		JobID:  *job.Id,
	}
	media.ID = metadata.ID

	err = repository.GetRepository().Insert(media)
	if err != nil {
		log.Println(err)
		err = tx.Rollback()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}

	for _, med := range mediainfo.Medias {
		for _, tra := range med.Tracks {
			track := &mediaModel.Track{
				MediaID:      media.ID,
				Type:         tra.Type,
				Duration:     tra.Duration,
				Width:        tra.Width,
				Height:       tra.Height,
				Format:       tra.Format,
				Encoded_Date: tra.Encoded_Date,
				VideoCount:   tra.VideoCount,
				DataSize:     tra.DataSize,
				FileSize:     tra.FileSize,
			}

			err = repository.GetRepository().Insert(track)
			if err != nil {
				log.Println(err)
				err = tx.Rollback()
				if err != nil {
					log.Fatal(err)
				}
				return nil
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
