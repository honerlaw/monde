package model

import "server/core/repository"

type Hashtag struct {
	repository.Model
	Tag        string
	MediaInfos []MediaInfo
}
