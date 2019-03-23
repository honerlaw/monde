package model

import (
	"github.com/jinzhu/gorm"
	"server/model/media"
)

type User struct {
	gorm.Model
	Username   string            `json:"username" gorm:"index;unique;not null"`
	Hash       string            `json:"hash" gorm:"not null"`
	MediaInfos []media.MediaInfo `gorm:"foreignkey:UserID"`
}

func (*User) Migrate(db *gorm.DB, migrate func(model interface{})) {
	migrate(&User{});
}
