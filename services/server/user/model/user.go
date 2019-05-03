package model

import (
	"services/server/core/repository"
)

type User struct {
	repository.Model
	Hash          string `json:"hash" column:"hash"`
}
