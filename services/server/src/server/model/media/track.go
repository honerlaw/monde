package media

import (
	"github.com/jinzhu/gorm"
	"encoding/xml"
)

type Track struct {
	gorm.Model
	MediaID      uint
	XMLName      xml.Name `xml:"track"`
	Type         string   `xml:"type,attr"`
	Duration     float64  `xml:"Duration"`
	Width        int64    `xml:"Width"`
	Height       int64    `xml:"Height"`
	Format       string   `xml:"Format"`
	Encoded_Date string   `xml:"Encoded_Date"`
	VideoCount   string   `xml:"VideoCount"`
	DataSize     int64    `xml:"DataSize"`
	FileSize     int64    `xml:"FileSize"`
}

func (*Track) Migrate(db *gorm.DB, migrate func(model interface{})) {
	migrate(&Track{})

	// have to do it this way because gorm does not add the foreign key for us
	db.Model(&Track{}).AddForeignKey("media_id", "media(id)", "CASCADE", "RESTRICT")
}
