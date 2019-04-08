package test

import (
	"github.com/joho/godotenv"
	"services/server/core/service/aws"
	"services/server/core/repository"
	"log"
	"os"
)

// NOTE: these can not be used in repository package for cyclic reasons

type TestModel struct {
	repository.Model
	TestField string `json:"test_field" column:"test_field"`
}

func Setup(rootPath string) {
	err := godotenv.Load(rootPath + "/.env")
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("DB_NAME", "vueon_test")
	os.Setenv("DB_MIGRATE_PATH", rootPath + "/migrations")

	// init aws session
	err = aws.InitSession()
	if err != nil {
		log.Fatal(err)
	}

	_, err = repository.GetRepository().Migrate().DB().Exec(`
CREATE TABLE IF NOT EXISTS ` + "`test_model`" + ` (
  ` + "`id`" + ` varchar(255) NOT NULL,
  ` + "`created_at`" + ` timestamp NULL DEFAULT NULL,
  ` + "`updated_at`" + ` timestamp NULL DEFAULT NULL,
  ` + "`deleted_at`" + ` timestamp NULL DEFAULT NULL,
  ` + "`test_field`" + ` varchar(255) NOT NULL,
  PRIMARY KEY (` + "`id`" + `),
  KEY ` + "`idx_media_deleted_at`" + ` (` + "`deleted_at`" + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`)

	if err != nil {
		log.Fatal(err)
	}
}

func Teardown() {
	_, err := repository.GetRepository().DB().Exec("DROP DATABASE " + os.Getenv("DB_NAME"))

	if err != nil {
		log.Fatal(err)
	}
}
