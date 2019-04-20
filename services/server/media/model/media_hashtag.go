package model

import "services/server/core/repository"

type MediaHashtag struct {
	repository.Model
	MediaID   string `json:"tag" column:"media_id"`
	HashtagID string `json:"tag" column:"hashtag_id"`
}
