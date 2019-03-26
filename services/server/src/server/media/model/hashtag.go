package model

import "github.com/jinzhu/gorm"

type Hashtag struct {
	gorm.Model
	Tag        string      `gorm:"unique_index"`
	MediaInfos []MediaInfo `gorm:"many2many:media_info_hashtag"`
}

func (*Hashtag) Migrate(db *gorm.DB, migrate func(model interface{})) {
	migrate(&Hashtag{})

	// @todo we need to add the foreign keys manually to the media_info_hashtag table
}
