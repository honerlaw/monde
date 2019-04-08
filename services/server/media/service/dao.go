package service

import (
	"services/server/media/model"
	"services/server/core/util"
	"log"
	"github.com/Masterminds/squirrel"
	"services/server/core/repository"
	"errors"
)

type MediaData struct {
	Media   model.Media
	Tracks []model.Track
	Tags   []model.Hashtag
}

func GetMediaData(selectPage *util.SelectPage) (*[]MediaData, error) {
	var medias []MediaData

	_, err := squirrel.Select("*").
		From("media_info mi").
		LeftJoin("media m ON mi.id = m.media_info_id").
		LeftJoin("track t ON m.id = t.media_id").
		Join("media_info_hashtag mih ON ht.media_info_id = mi.id").
		Join("hashtag h ON h.id = mih.hashtag_id").
		RunWith(repository.GetRepository().DB()).
		Query()

	if err != nil {
		log.Print("failed to get media info for user", err)
		return nil, errors.New("failed to find media information")
	}

	return &medias, nil
}

func GetMediaDataByUserId(userId string, selectPage *util.SelectPage) (*[]MediaData, error) {
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

func Save(media *MediaData) (error) {
	tx, err := repository.GetRepository().DB().Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
