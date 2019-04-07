package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"os"
	"fmt"
	"encoding/json"
	"server/core/service/aws"
	"sync"
	"strconv"
	"database/sql"
	"log"
	"time"
)

var dbOnce sync.Once
var repoInstance *Repository

type Model struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type DBCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DBInfo struct {
	creds    *DBCredentials
	endpoint string
	dbname   string
}

type Repository struct {
	url string
	db  *sql.DB
}

func GetRepository() (*Repository) {
	dbOnce.Do(func() {
		info := getDBInfo()

		// create the db
		url := generateDbUrl(info.creds, info.endpoint, nil)
		err := createDatabaseIfNotExists(url, info.dbname)
		if err != nil {
			log.Fatal(err)
		}

		// get the actual db connection
		url = generateDbUrl(info.creds, info.endpoint, &info.dbname)
		db, err := sql.Open("mysql", url)
		if err != nil {
			log.Fatal(err)
		}

		repoInstance = &Repository{
			url: url,
			db:  db,
		}
	})
	return repoInstance
}

func (repo *Repository) DB() (*sql.DB) {
	return repo.db
}

func (repo *Repository) Migrate() (*Repository) {
	shouldMigrate, err := strconv.ParseBool(os.Getenv("DB_MIGRATE"))
	if !shouldMigrate || err != nil {
		return repo
	}

	driver, _ := mysql.WithInstance(repo.db, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "mysql", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}

	return repo
}

func getDBInfo() (*DBInfo) {
	// get the database cluster information
	cluster, err := aws.GetRDSService().GetCluster()
	if err != nil {
		log.Fatal(err)
	}

	// get the database credentials secret
	secret, err := aws.GetSMService().GetSecret(os.Getenv("DB_SECRET_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	// secret is stored as json so unmarshal it
	var creds DBCredentials
	json.Unmarshal([]byte(secret.Value), &creds)

	return &DBInfo{
		creds:    &creds,
		endpoint: cluster.Endpoint,
		dbname:   os.Getenv("DB_NAME"),
	}
}

func generateDbUrl(creds *DBCredentials, endpoint string, dbname *string) (string) {
	if dbname != nil {
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true", creds.Username, creds.Password, endpoint, *dbname)
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/?parseTime=true&multiStatements=true", creds.Username, creds.Password, endpoint)
}

func createDatabaseIfNotExists(dbUrl string, dbname string) (error) {
	// attempt to create the database first
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close();
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
