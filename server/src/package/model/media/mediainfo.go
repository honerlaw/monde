package media

import (
	"github.com/jinzhu/gorm"
	"encoding/xml"
)

type MediaInfo struct {
	gorm.Model
	XMLName       xml.Name `xml:"MediaInfo"`
	Medias        []Media `xml:"media" gorm:"foreignkey:MediaInfoID"`
	UserID        uint
	JobID         string
	VideoID       string
}

func (*MediaInfo) Migrate(db *gorm.DB, migrate func(model interface{})) {
	migrate(&MediaInfo{})

	// have to do it this way because gorm does not add the foreign key for us
	db.Model(&MediaInfo{}).AddForeignKey("user_id", "user(id)", "CASCADE", "RESTRICT")
}
