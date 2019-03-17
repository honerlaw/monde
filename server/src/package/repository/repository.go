package repository

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"package/model"
	"package/model/media"
	"reflect"
	"log"
	"os"
	"strings"
	"fmt"
)

var DB *gorm.DB

func Connect(url string) *gorm.DB {
	noDbUrl := strings.Replace(url, "/" + os.Getenv("DB_NAME"), "/", -1)

	// attempt to create the database first
	tempDb, err := gorm.Open("mysql", noDbUrl)
	if err != nil {
		panic(err)
	}
	tempDb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", os.Getenv("DB_NAME")))
	if tempDb.Error != nil {
		panic(tempDb.Error)
	}
	tempDb.Close();

	// then get the normal database connection
	DB, err = gorm.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	DB.SingularTable(true)

	return DB
}

func Migrate() {
	(&model.User{}).Migrate(DB, migrateModel)
	(&media.MediaInfo{}).Migrate(DB, migrateModel)
	(&media.Media{}).Migrate(DB, migrateModel)
	(&media.Track{}).Migrate(DB, migrateModel)
}

func migrateModel(model interface{}) {
	modelType := reflect.TypeOf(model)
	log.Printf("Auto Migrating %s\n", modelType)

	db := DB.AutoMigrate(model)
	if db.Error != nil {
		log.Printf("%s failed with error %s\n", modelType, db.Error)
		panic(db.Error)
	}
	log.Printf("%s was successfully migrated\n", modelType)
}
