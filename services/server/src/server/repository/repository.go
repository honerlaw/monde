package repository

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"server/model"
	"server/model/media"
	"reflect"
	"log"
	"os"
	"fmt"
	"server/service/aws"
	"encoding/json"
)

var DB *gorm.DB

type DBCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Connect() *gorm.DB {
	// get the database cluster information
	cluster, err := aws.GetRDSService().GetCluster()
	if err != nil {
		panic(err)
	}

	// get the database credentials secret
	secret, err := aws.GetSMService().GetSecret(os.Getenv("DB_SECRET_NAME"))
	if err != nil {
		panic(err)
	}

	// secret is stored as json so unmarshal it
	var creds DBCredentials
	json.Unmarshal([]byte(secret.Value), &creds)

	// create the db
	err = createDatabaseIfNotExists(&creds, cluster.Endpoint)
	if err != nil {
		panic(err)
	}

	// get the actual db connection
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", creds.Username, creds.Password, cluster.Endpoint, os.Getenv("DB_NAME"))
	DB, err = gorm.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	// all tables are named singluar, e.g. user instead of users
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

func createDatabaseIfNotExists(creds *DBCredentials, endpoint string) (error) {

	// get the db url without the database name
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s)/?parseTime=true", creds.Username, creds.Password, endpoint)

	// attempt to create the database first
	tempDb, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
	}
	tempDb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", os.Getenv("DB_NAME")))
	if tempDb.Error != nil {
		panic(tempDb.Error)
	}
	tempDb.Close();
	return nil
}