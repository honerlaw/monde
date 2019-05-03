package model

import "services/server/core/repository"

type Contact struct {
	repository.Model
	UserID   string `json:"user_id" column:"user_id"`
	Contact  string `json:"contact" column:"contact"`
	Type     string `json:"type" column:"type"`
	Code     string `json:"code" column:"code"`
	Verified bool   `json:"verified" column:"verified"`
}
