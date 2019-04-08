package model

import (
	"time"
	"strings"
	"services/server/core/repository"
)

type Media struct {
	repository.Model
	UserID        string
	JobID         string
	Title         string
	Description   string
	Published     bool
	PublishedDate time.Time
}

func (info *Media) CanPublish() (bool) {
	return len(strings.TrimSpace(info.Description)) > 0
}
