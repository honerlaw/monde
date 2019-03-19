package media

import (
	"package/model/media"
	"package/repository"
	"log"
)

func GetMediaInfoByUserId(userId uint) (*[]media.MediaInfo, error) {
	var infos []media.MediaInfo

	repository.DB.Where(media.MediaInfo{UserID: userId}).Order("created_at DESC").Find(&infos)

	if repository.DB.Error != nil {
		log.Print("failed to get media info for user", repository.DB.Error)
		return nil, repository.DB.Error
	}

	return &infos, nil
}
