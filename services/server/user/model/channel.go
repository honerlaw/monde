package model

import (
	"services/server/core/repository"
)

type Channel struct {
	repository.Model
	UserID string `json:"user_id" column:"user_id"`
	Title string `json:"title" column:"title"`
	Slug string `json:"slug" column:"slug"`
}
