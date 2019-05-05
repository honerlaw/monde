package repository

import (
	"services/server/core/repository"
	"reflect"
	"services/server/user/model"
	"github.com/Masterminds/squirrel"
	"github.com/labstack/gommon/log"
)

type ContactRepository struct {
	repo *repository.Repository
}

func NewContactRepository(repo *repository.Repository) (*ContactRepository) {
	return &ContactRepository{
		repo: repo,
	}
}

func (repo *ContactRepository) tableName() (string) {
	modelValue := reflect.Indirect(reflect.ValueOf(&model.Contact{}))
	modelType := modelValue.Type()
	return repo.repo.Table(modelType)
}

func (repo *ContactRepository) FindByUserID(userID string) ([]model.Contact, error) {

	rows, err := squirrel.Select("*").
		From(repo.tableName()).
		Where(squirrel.Eq{"user_id": userID}).
		RunWith(repo.repo.DB()).
		Query()

	if err != nil {
		log.Print(err)
		return nil, err
	}

	parsed := repo.repo.Parse(&model.Contact{}, rows)
	models := make([]model.Contact, len(parsed))
	for i := 0; i < len(parsed); i++ {
		models[i] = parsed[i].(model.Contact)
	}

	return models, nil
}

func (repo *ContactRepository) FindByContact(contact string) (*model.Contact, error) {

	// we assume that contact has been normalized by the time it has gotten here
	rows, err := squirrel.Select("*").
		From(repo.tableName()).
		Where(squirrel.Eq{"contact": contact}).
		RunWith(repo.repo.DB()).
		Query()

	if err != nil {
		log.Print(err)
		return nil, err
	}

	parsed := repo.repo.Parse(&model.Contact{}, rows)
	if len(parsed) != 1 {
		return nil, nil
	}
	contactModel := parsed[0].(model.Contact)

	return &contactModel, nil
}

func (repo *ContactRepository) Save(contact *model.Contact) (error) {
	return repo.repo.Save(contact)
}
