package repository

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"server/model"
	"reflect"
	"log"
	"os"
	"fmt"
	"server/service/aws"
	mediaModel "server/media/model"
	"encoding/json"
)

var DB *gorm.DB

type DBCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DBInfo struct {
	creds *DBCredentials
	endpoint string
	dbname string
}

func Connect() *gorm.DB {

	info := getDBInfo()

	// create the db
	url := generateDbUrl(info.creds, info.endpoint, nil)
	err := createDatabaseIfNotExists(url, info.dbname)
	if err != nil {
		panic(err)
	}

	// get the actual db connection
	url = generateDbUrl(info.creds, info.endpoint, &info.dbname)
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
	(&mediaModel.MediaInfo{}).Migrate(DB, migrateModel)
	(&mediaModel.Media{}).Migrate(DB, migrateModel)
	(&mediaModel.Track{}).Migrate(DB, migrateModel)
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

func getDBInfo() (*DBInfo) {
	if os.Getenv("ENV") == "local" {
		return &DBInfo{
			creds: &DBCredentials{
				Username: os.Getenv("DB_USER_LOCAL"),
				Password: os.Getenv("DB_PASS_LOCAL"),
			},
			endpoint: os.Getenv("DB_ENDPOINT_LOCAL"),
			dbname: os.Getenv("DB_NAME"),
		}
	}

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

	return &DBInfo{
		creds: &creds,
		endpoint: cluster.Endpoint,
		dbname: os.Getenv("DB_NAME"),
	}
}

func generateDbUrl(creds *DBCredentials, endpoint string, dbname *string) (string) {
	if dbname != nil {
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", creds.Username, creds.Password, endpoint, *dbname)
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/?parseTime=true", creds.Username, creds.Password, endpoint)
}

func createDatabaseIfNotExists(dbUrl string, dbname string) (error) {
	// attempt to create the database first
	tempDb, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
	}
	tempDb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname))
	if tempDb.Error != nil {
		panic(tempDb.Error)
	}
	tempDb.Close();
	return nil
}