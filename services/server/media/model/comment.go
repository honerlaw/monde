package model

import "services/server/core/repository"

type Comment struct {
	repository.Model
	MediaID         string `json:"media_id" column:"media_id"`
	ParentCommentID string `json:"parent_comment_id" column:"parent_comment_id"`
	Comment         string `json:"comment" column:"comment"`
}
