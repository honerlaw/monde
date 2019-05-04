package model

import (
	"services/server/core/repository"
)

type Address struct {
	repository.Model
	UserID  string `json:"user_id" column:"user_id"`
	Type    string `json:"type" column:"type"`
	LineOne string `json:"line_one" column:"line_one"`
	LineTwo string `json:"line_two" column:"line_two"`
	City    string `json:"city" column:"city"`
	State   string `json:"state" column:"state"`
	ZipCode string `json:"zip_code" column:"zip_code"`
	Country string `json:"country" column:"country"`
}
