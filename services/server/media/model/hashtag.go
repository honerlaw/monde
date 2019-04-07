package model

import "services/server/core/repository"

type Hashtag struct {
	repository.Model
	Tag        string
	MediaInfos []MediaInfo
}
