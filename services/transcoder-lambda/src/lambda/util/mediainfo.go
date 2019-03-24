package util

import (
	"os/exec"
	"lambda/aws"
	"encoding/xml"
	"errors"
	"server/media/model"
	"log"
	"fmt"
	serverAWS "server/service/aws"
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

// this monstrosity should insert all of the media info properly
// We can't really use transactions because of the wonderful aurora data api limitations...
// We shouldn't need to worry about escaping things because we generated all the data...
// HOWEVER, we should probably still try and escape crap just in case mediainfo gets some weird results back
// @todo escape things to prevent sql injections / unintended side effects
func Insert(mediainfo *model.MediaInfo) (error) {
	service := aws.GetRDSDService()

	_, err := service.ExecuteSQL(fmt.Sprintf(
		`INSERT INTO media_info (created_at, user_id, job_id, video_id) VALUES (NOW(), %d, '%s', '%s')`,
		mediainfo.UserID, mediainfo.JobID, mediainfo.VideoID,
	))
	if err != nil {
		return err
	}

	resp, err := service.ExecuteSQL(fmt.Sprintf(
		"SELECT id FROM media_info WHERE user_id = %d AND job_id = '%s' and video_id = '%s'",
		mediainfo.UserID, mediainfo.JobID, mediainfo.VideoID,
	))
	if err != nil {
		return err
	}

	mediaInfoId := *resp.SqlStatementResults[0].ResultFrame.Records[0].Values[0].IntValue;
	for _, mediaInfoMedia := range mediainfo.Medias {
		_, err = service.ExecuteSQL(fmt.Sprintf(
			`INSERT INTO media (created_at, media_info_id) VALUES (NOW(), %d)`,
			mediaInfoId,
		))
		if err != nil {
			return err
		}

		resp, err := service.ExecuteSQL(fmt.Sprintf(
			"SELECT id FROM media WHERE media_info_id = %d ORDER BY created_at DESC LIMIT 1",
			mediaInfoId,
		))
		if err != nil {
			return err
		}
		mediaId := *resp.SqlStatementResults[0].ResultFrame.Records[0].Values[0].IntValue

		for _, track := range mediaInfoMedia.Tracks {
			_, err = service.ExecuteSQL(fmt.Sprintf(
				`INSERT INTO track (created_at, media_id, type, duration, width, height, format, encoded_date, video_count, data_size, file_size) VALUES (NOW(), %d, '%s', %f, %d, %d, '%s', '%s', '%s', %d, %d)`,
				mediaId, track.Type, track.Duration, track.Width, track.Height, track.Format, track.Encoded_Date, track.VideoCount, track.DataSize, track.FileSize,
			))
			if err != nil {
				return err
			}
		}
	}

	return nil
}