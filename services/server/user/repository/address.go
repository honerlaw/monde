package repository

import (
	"services/server/core/repository"
	"services/server/user/model"
	"github.com/Masterminds/squirrel"
	"reflect"
	"github.com/labstack/gommon/log"
)

type AddressRepository struct {
	repo *repository.Repository
}

func NewAddressRepository(repo *repository.Repository) (*AddressRepository) {
	return &AddressRepository{
		repo: repo,
	}
}

func (repo *AddressRepository) tableName() (string) {
	modelValue := reflect.Indirect(reflect.ValueOf(&model.Channel{}))
	modelType := modelValue.Type()
	return repo.repo.Table(modelType)
}

func (repo *AddressRepository) FindByUserID(userID  string) ([]model.Address, error) {
	rows, err := squirrel.Select("*").
		From(repo.tableName()).
		Where(squirrel.Eq{"user_id": userID}).
		RunWith(repo.repo.DB()).
		Query()

	if err != nil {
		log.Print(err)
		return nil, err
	}

	parsed := repo.repo.Parse(&model.Address{}, rows)
	// @todo convert everything

	return nil, nil
}

func (repo *AddressRepository) Save(model *model.Address) (error) {
	return repo.repo.Save(model)
}
