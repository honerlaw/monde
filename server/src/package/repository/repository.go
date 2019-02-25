package repository

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"package/model"
)

var DB *gorm.DB

func Connect(url string) *gorm.DB {
	var err error
	DB, err = gorm.Open("mysql", url)
	if err != nil {
		panic(err)
	}
	return DB
}

func Migrate() {
	DB.AutoMigrate(&model.User{})
}