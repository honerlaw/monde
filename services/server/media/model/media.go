package model

import (
	"time"
	"strings"
	"services/server/core/repository"
)

type Media struct {
	repository.Model
	ChannelID     string    `json:"channel_id" column:"channel_id"`
	JobID         string    `json:"job_id" column:"job_id"`
	Title         string    `json:"title" column:"title"`
	Description   string    `json:"description" column:"description"`
	Published     bool      `json:"published" column:"published"`
	PublishedDate time.Time `json:"published_date" column:"published_date"`
}

func (info *Media) CanPublish() (bool) {
	return len(strings.TrimSpace(info.Description)) > 0
}
