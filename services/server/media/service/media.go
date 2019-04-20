package service

import (
	"time"
	"strings"
	"regexp"
	"services/server/media/model"
	"errors"
	"log"
	"services/server/core/repository"
)

type UpdateRequest struct {
	VideoID     string `form:"video_id" binding:"required"`
	Title       string `form:"title"`
	Description string `form:"description"`
	Hashtags    string `form:"hashtags"`
}

type PublishRequest struct {
	VideoID string `form:"video_id" binding:"required"`
}

func Update(req UpdateRequest) (error) {
	tx, err := repository.GetRepository().DB().Begin()

	if err != nil {
		log.Print("failed to start transaction", err)
		return err
	}

	data, err := GetMediaDataByVideoID(req.VideoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	data.Media.Title = req.Title
	data.Media.Description = req.Description

	// if they remove the description, unpublish the video
	if len(strings.TrimSpace(data.Media.Description)) == 0 {
		data.Media.Published = false;
	}

	var hashtags []model.Hashtag
	regex, _ := regexp.Compile("^#\\w+")
	tags := strings.Split(req.Hashtags, " ")
	for _, tag := range tags {
		trimmedTag := strings.TrimSpace(tag)
		if regex.MatchString(trimmedTag) && len(trimmedTag) < 75 && len(trimmedTag) > 1 {
			// fetch the existing hashtag if it exists
			hashtag, err := GetHashTag(trimmedTag)
			if hashtag == nil || err != nil {
				hashtag = &model.Hashtag{
					Tag: trimmedTag,
				}
			}

			hashtags = append(hashtags, *hashtag)
		}
	}
	data.Tags = hashtags

	err = Save(data)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func TogglePublish(req PublishRequest) (error) {
	data, err := GetMediaDataByVideoID(req.VideoID)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(data.Media.Description)) == 0 {
		return errors.New("a description is required to publish videos")
	}

	data.Media.Published = !data.Media.Published
	if data.Media.Published {
		data.Media.PublishedDate = time.Now()
	}

	err = Save(data)
	if err != nil {
		return err
	}

	return nil
}
