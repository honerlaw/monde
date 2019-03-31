package model

import (
	"github.com/jinzhu/gorm"
	"encoding/xml"
	"time"
	"strings"
)

type MediaInfo struct {
	gorm.Model
	XMLName       xml.Name  `xml:"MediaInfo"`
	Medias        []Media   `xml:"media" gorm:"foreignkey:MediaInfoID"`
	Hashtags      []Hashtag `gorm:"many2many:media_info_hashtag"`
	UserID        uint
	JobID         string
	VideoID       string
	Title         string `gorm:"type:tinytext"`
	Description   string `gorm:"type:text"`
	Published     bool
	PublishedDate time.Time
}

func (info *MediaInfo) CanPublish() (bool) {
	return len(strings.TrimSpace(info.Description)) > 0
}

func (*MediaInfo) Migrate(db *gorm.DB, migrate func(model interface{})) {
	migrate(&MediaInfo{})

	// have to do it this way because gorm does not add the foreign key for us
	db.Model(&MediaInfo{}).AddForeignKey("user_id", "user(id)", "CASCADE", "RESTRICT")
}
