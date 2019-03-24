package service

import (
	"server/repository"
	"log"
	"server/media/model"
)

type UpdateRequest struct {
	VideoID string `form:"video_id" binding:"required"`
	Title string `form:"title"`
	Description string `form:"description"`
	Hashtags string `form:"hashtags"`
}

func GetMediaInfoByUserId(userId uint) (*[]model.MediaInfo, error) {
	var infos []model.MediaInfo

	repository.DB.Where(model.MediaInfo{UserID: userId}).Order("created_at DESC").Find(&infos)

	if repository.DB.Error != nil {
		log.Print("failed to get media info for user", repository.DB.Error)
		return nil, repository.DB.Error
	}

	return &infos, nil
}

func Update(req UpdateRequest) (error) {
	return nil
}


