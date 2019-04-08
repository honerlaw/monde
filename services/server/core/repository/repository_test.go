package repository

import (
	"testing"
	"os"
	"log"
	"github.com/joho/godotenv"
	"services/server/core/service/aws"
)

type testModel struct {
	Model
	TestField string `json:"test_field" column:"test_field"`
}

func setup() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	// init aws session
	err = aws.InitSession()
	if err != nil {
		log.Fatal(err)
	}

	_, err = GetRepository().DB().Exec(`
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

func teardown() {
	_, err := GetRepository().DB().Exec("DROP TABLE `test_model`")

	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func TestInsert(t *testing.T) {
	temp := testModel{
		TestField: "some_value",
	}

	err := GetRepository().Insert(&temp)
	if err != nil {
		t.Error(err)
	}

	found, err := GetRepository().FindByID(temp.ID, &temp)
	if !found || err != nil {
		t.Error(found, err)
	}
}

func TestUpdate(t *testing.T) {
	temp := testModel{
		TestField: "some_value",
	}

	// do the initial insert
	err := GetRepository().Insert(&temp)
	if err != nil {
		t.Error(err)
	}

	// change value for update
	temp.TestField = "testing"

	// trigger update
	err = GetRepository().Update(&temp)
	if err != nil {
		t.Error(err)
	}

	// verify that the data was updated and exists
	found, err := GetRepository().FindByID(temp.ID, &temp)
	if !found || err != nil {
		t.Error(found, err)
	}

	// verify that the update actual persisted
	if temp.TestField != "testing" {
		t.Error("update failed")
	}
}

func TestSave(t *testing.T) {
	temp := testModel{
		TestField: "some_value",
	}

	err := GetRepository().Save(&temp)
	if err != nil {
		t.Error(err)
	}

	// verify that the data was updated and exists
	found, err := GetRepository().FindByID(temp.ID, &temp)
	if !found || err != nil || temp.TestField != "some_value" {
		t.Error(found, err, temp.TestField)
	}

	temp.TestField = "existing"
	err = GetRepository().Save(&temp)
	if err != nil {
		t.Error(err)
	}

	// verify that the data was updated and exists
	found, err = GetRepository().FindByID(temp.ID, &temp)
	if !found || err != nil || temp.TestField != "existing" {
		t.Error(found, err, temp.TestField)
	}
}