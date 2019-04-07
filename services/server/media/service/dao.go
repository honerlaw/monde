package service

import (
	"services/server/media/model"
	"services/server/core/util"
)

type MediaData struct {
	Info   model.MediaInfo
	Tags   []model.Hashtag
	Medias []model.Media
	Tracks []model.Track
}

func GetMediaData(selectPage *util.SelectPage) (*[]MediaData, error) {
	return nil, nil
	/*var infos []model.MediaInfo

	// @todo optimize this? makes 3 queries...
	repository.DB.Where(&model.MediaInfo{
		Published: true,
	}).Order("created_at DESC").
		Offset(selectPage.Page).
		Limit(selectPage.Count).
		Preload("Hashtags").
		Preload("Medias").
		Preload("Medias.Tracks").
		Find(&infos)

	if repository.DB.Error != nil {
		log.Print("failed to get media info for user", repository.DB.Error)
		return nil, errors.New("failed to find media information")
	}

	return &infos, nil*/
}

func GetMediaDataByUserId(userId uint, selectPage *util.SelectPage) (*[]MediaData, error) {
	return nil, nil/*
	var infos []model.MediaInfo

	offset := selectPage.Page * selectPage.Count

	// @todo optimize this? makes 3 queries...
	repository.DB.Where(model.MediaInfo{UserID: userId}).
		Order("created_at DESC").
		Offset(offset).Limit(selectPage.Count).
		Preload("Hashtags").
		Preload("Medias").
		Preload("Medias.Tracks").
		Find(&infos)

	if repository.DB.Error != nil {
		log.Print("failed to get media info for user", repository.DB.Error)
		return nil, errors.New("failed to find media information")
	}

	return &infos, nil*/
}

func GetMediaDataByVideoID(videoId string) (*MediaData, error) {
	return nil, nil /*
	var info model.MediaInfo

	repository.DB.
		Where(model.MediaInfo{VideoID: videoId}).
		Order("created_at DESC").
		Preload("Hashtags").
		Preload("Medias").
		Preload("Medias.Tracks").
		First(&info)

	if repository.DB.Error != nil {
		log.Print("failed to get media info for user", repository.DB.Error)
		return nil, errors.New("failed to find media information")
	}

	return &info, nil*/
}

func GetHashTag(tag string) (*model.Hashtag, error) {
	return nil, nil/*
	var hashtag model.Hashtag

	repository.DB.Where(model.Hashtag{Tag: tag}).First(&hashtag)
	if repository.DB.Error != nil {
		log.Printf("failed to find tag for tag: %s, error: %s", tag, repository.DB.Error.Error())
		return nil, errors.New("failed to find hashtag")
	}

	// we didn't find the tag
	if hashtag.Tag == "" {
		return nil, nil
	}

	return &hashtag, nil*/
}

func Save(data *MediaData) (error) {
	/*repository.DB.Save(model)
	if repository.DB.Error != nil {
		log.Print("failed to save model", repository.DB.Error)
		return errors.New("failed to save information")
	}*/
	return nil
}
