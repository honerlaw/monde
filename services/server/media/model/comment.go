package model

import "services/server/core/repository"

type Comment struct {
	repository.Model
	ParentCommentID string `json:"parent_comment_id" column:"parent_comment_id"`
	MediaID         string `json:"media_id" column:"media_id"`
	UserID          string `json:"user_id" column:"user_id"`
	Comment         string `json:"comment" column:"comment"`
}
