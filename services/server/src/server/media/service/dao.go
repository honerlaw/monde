package service

import (
	"server/media/model"
	"server/repository"
	"log"
	"errors"
)

type SelectPage struct {
	Page int
	Count int
}

func GetMediaInfoByUserId(userId uint, selectPage *SelectPage) (*[]model.MediaInfo, error) {
	var infos []model.MediaInfo

	offset := selectPage.Page * selectPage.Count

	repository.DB.Where(model.MediaInfo{UserID: userId}).Order("created_at DESC").Offset(offset).Limit(selectPage.Count).Find(&infos)

	if repository.DB.Error != nil {
		log.Print("failed to get media info for user", repository.DB.Error)
		return nil, errors.New("failed to find media information")
	}

	return &infos, nil
}

func GetMediaInfoByVideoID(videoId string) (*model.MediaInfo, error) {
	var info model.MediaInfo

	repository.DB.Where(model.MediaInfo{VideoID: videoId}).Order("created_at DESC").First(&info)

	if repository.DB.Error != nil {
		log.Print("failed to get media info for user", repository.DB.Error)
		return nil, errors.New("failed to find media information")
	}

	return &info, nil
}

func GetHashTag(tag string) (*model.Hashtag, error) {
	var hashtag model.Hashtag

	repository.DB.Where(model.Hashtag{ Tag: tag }).First(&hashtag)
	if repository.DB.Error != nil {
		log.Printf("failed to find tag for tag: %s, error: %s", tag, repository.DB.Error.Error())
		return nil, errors.New("failed to find hashtag")
	}

	// we didn't find the tag
	if hashtag.Tag == "" {
		return nil, nil
	}

	return &hashtag, nil
}

func Save(model interface{}) (error) {
	repository.DB.Save(model)
	if repository.DB.Error != nil {
		log.Print("failed to save model", repository.DB.Error)
		return errors.New("failed to save information")
	}
	return nil
}
