package model

import (
	"github.com/jinzhu/gorm"
	"encoding/xml"
)

type Media struct {
	gorm.Model
	XMLName     xml.Name `xml:"media"`
	Tracks      []Track  `xml:"track" gorm:"foreignkey:MediaID"`
	MediaInfoID uint
}

func (*Media) Migrate(db *gorm.DB, migrate func(model interface{})) {
	migrate(&Media{})

	// have to do it this way because gorm does not add the foreign key for us
	db.Model(&Media{}).AddForeignKey("media_info_id", "media_info(id)", "CASCADE", "RESTRICT")
}
