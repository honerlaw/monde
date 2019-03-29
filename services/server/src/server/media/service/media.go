package service

import (
	"time"
	"strings"
	"regexp"
	"server/media/model"
	"errors"
	"log"
	"server/core/repository"
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
	tx := repository.DB.Begin()

	info, err := GetMediaInfoByVideoID(req.VideoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	info.Title = req.Title
	info.Description = req.Description

	err = Save(info)
	if err != nil {
		tx.Rollback()
		return err
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

	assoc := tx.Model(info).Association("Hashtags").Replace(hashtags)
	if assoc.Error != nil {
		tx.Rollback()
		log.Print(assoc.Error)
		return errors.New("failed to update media")
	}

	tx.Commit()

	return nil
}

func TogglePublish(req PublishRequest) (error) {
	info, err := GetMediaInfoByVideoID(req.VideoID)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(info.Description)) == 0 {
		return errors.New("a description is required to publish videos")
	}

	info.Published = !info.Published
	if info.Published {
		info.PublishedDate = time.Now()
	}

	err = Save(info)
	if err != nil {
		return err
	}

	return nil
}

