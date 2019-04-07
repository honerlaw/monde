package model

import (
	"services/server/core/repository"
)

// @todo potentally use reflection for this info? then we can have some generic ways for save / parsing
var UserColumns = []string{"id", "created_at", "updated_at", "deleted_at", "username", "hash"}

type User struct {
	repository.Model
	Username string `json:"username"`
	Hash     string `json:"hash"`
}
