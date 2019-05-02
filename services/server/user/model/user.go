package model

import (
	"services/server/core/repository"
)

type User struct {
	repository.Model
	Email         string `json:"email" column:"email"`
	Username      string `json:"username" column:"username"`
	Hash          string `json:"hash" column:"hash"`
	VerifiedEmail bool   `json:"verified_email" column:"verified_email"`
}
