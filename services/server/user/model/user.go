package model

import (
	"services/server/core/repository"
)

type User struct {
	repository.Model
	Username string `json:"username" column:"username"`
	Hash     string `json:"hash" column:"hash"`
}
