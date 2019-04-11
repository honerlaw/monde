package repository

import (
	"testing"
	"os"
	"github.com/joho/godotenv"
	"services/server/core/service/aws"
	"log"
	"github.com/Masterminds/squirrel"
	"reflect"
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
	os.Setenv("DB_NAME", "vueon_test")
	os.Setenv("DB_MIGRATE_PATH",  "../../migrations")

	// init aws session
	err = aws.InitSession()
	if err != nil {
		log.Fatal(err)
	}

	_, err = GetRepository().Migrate().DB().Exec(`
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
	_, err := GetRepository().DB().Exec("DROP DATABASE " + os.Getenv("DB_NAME"))

	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	setup()

	// make sure to teardown on panic
	defer func() {
		if r := recover(); r != nil {
			teardown()
		}
	}()

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

	found, err := GetRepository().FindByID(temp.ID, &testModel{})
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

func TestDelete(t *testing.T) {
	temp := testModel{
		TestField: "some_value",
	}

	err := GetRepository().Save(&temp)
	if err != nil {
		t.Error(err)
	}

	// verify that the data was updated and exists
	found, err := GetRepository().FindByID(temp.ID, &temp)
	if !found || err != nil {
		t.Error(found, err)
	}

	err = GetRepository().Delete(&temp)
	if err != nil {
		t.Error(err)
	}

	found, err = GetRepository().FindByID(temp.ID, &temp)
	if found || err != nil {
		t.Error(found, err)
	}
}

func TestParse(t *testing.T) {
	temp := testModel{
		TestField: "some_value",
	}

	err := GetRepository().Insert(&temp)
	if err != nil {
		t.Error(err)
	}

	table := GetRepository().Table(reflect.TypeOf(&temp).Elem())

	rows, err := squirrel.Select("*").
		From(table).
		Where(squirrel.Eq{"id": temp.ID}).
		RunWith(GetRepository().DB()).
		Query()

	if err != nil {
		t.Error(err)
	}

	models := GetRepository().Parse(&temp, rows)

	if len(models) == 0 {
		t.Error("one model should be found")
	}

	model := models[0].(testModel)

	if len(model.ID) == 0 {
		t.Error("Failed to properly parse and store information into new object")
	}
}
