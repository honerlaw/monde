package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"os"
	"fmt"
	"encoding/json"
	"services/server/core/service/aws"
	"sync"
	"strconv"
	"database/sql"
	"log"
	"time"
	"reflect"
	"github.com/Masterminds/squirrel"
	"services/server/core/util"
	"github.com/satori/go.uuid"
	"strings"
)

var dbOnce sync.Once
var repoInstance *Repository

type Model struct {
	ID        string    `json:"id" column:"id"`
	CreatedAt time.Time `json:"created_at" column:"created_at"`
	UpdatedAt time.Time `json:"updated_at" column:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" column:"deleted_at"`
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

func (repo *Repository) Parse(modelType reflect.Type, rows *sql.Rows) ([]interface{}) {
	defer rows.Close()

	models := make([]interface{}, 0)

	for rows.Next() {
		newModel := reflect.Zero(modelType).Interface()
		modelValue := reflect.ValueOf(newModel)
		pointers := repo.Values(modelValue)

		for i := 0; i < len(pointers); i++ {
			pointer := pointers[i]
			if reflect.TypeOf(pointer).Kind() != reflect.Ptr {
				pointers[i] = &pointer;
			}
		}

		// set the values
		err := rows.Scan(pointers...)
		if err != nil {
			log.Fatal(err)
		}

		models = append(models, newModel)
	}

	err := rows.Err()
	if err != nil {
		log.Print(err)
		return models
	}

	return models;
}

func (repo *Repository) FindByID(id string, model interface{}) (bool, error) {
	modelType := reflect.TypeOf(model).Elem()
	table := repo.Table(modelType)

	rows, err := squirrel.Select("*").
		From(table).
		Where(squirrel.Eq{"id": id}).
		RunWith(repo.db).
		Query()

	if err != nil {
		log.Print(err)
		return false, err
	}

	models := repo.Parse(modelType, rows)
	if len(models) > 0 {
		m := models[0]
		model = m
		return true, nil
	}
	return false, nil
}

func (repo *Repository) Save(model interface{}) (error) {
	modelValue := reflect.ValueOf(model).Elem()
	id := modelValue.FieldByName("ID").String()

	found, err := repo.FindByID(id, model)
	if err != nil {
		return err
	}

	if found {
		return repo.Update(model)
	}

	return repo.Insert(model)
}

func (repo *Repository) Update(model interface{}) (error) {
	modelType := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()

	modelValue.FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now()))
	id := modelValue.FieldByName("ID").String()

	table := repo.Table(modelType)
	columns := repo.Columns(modelType)
	values := repo.Values(modelValue)

	update := squirrel.Update(table)
	for index, column := range columns {
		update = update.Set(column, values[index])
	}
	update = update.Where(squirrel.Eq{"id": id})

	_, err := update.RunWith(repo.db).Query()
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (repo *Repository) Insert(model interface{}) (error) {
	modelType := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()

	id := uuid.NewV4().String()
	t := time.Now()
	modelValue.FieldByName("ID").SetString(id)
	modelValue.FieldByName("CreatedAt").Set(reflect.ValueOf(t))
	modelValue.FieldByName("UpdatedAt").Set(reflect.ValueOf(t))

	table := repo.Table(modelType)
	columns := repo.Columns(modelType)
	values := repo.Values(modelValue)

	_, err := squirrel.Insert(table).
		Columns(columns...).
		Values(values...).
		RunWith(repo.db).
		Query()

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (repo *Repository) Columns(modelType reflect.Type) ([]string) {
	columns := []string{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		if field.Type.Kind() == reflect.Struct && strings.Contains(field.Type.String(), "repository.Model") {
			columns = append(columns, repo.Columns(field.Type)...)
		} else {
			columns = append(columns, field.Tag.Get("column"))
		}
	}

	return columns
}

func (repo *Repository) Values(modelValue reflect.Value) ([]interface{}) {
	values := make([]interface{}, 0)

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		if field.Kind() == reflect.Struct && strings.Contains(field.String(), "repository.Model") {
			values = append(values, repo.Values(field)...)
		} else {
			values = append(values, field.Interface())
		}
	}

	return values
}

func (repo *Repository) Table(modelType reflect.Type) (string) {
	return util.ToSnakeCase(modelType.Name())
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
