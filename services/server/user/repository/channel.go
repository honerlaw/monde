package repository

import (
	"services/server/user/model"
	"services/server/core/repository"
	"github.com/labstack/gommon/log"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"reflect"
)

type ChannelRepository struct {
	repo *repository.Repository
}

func NewChannelRepository(repo *repository.Repository) (*ChannelRepository) {
	return &ChannelRepository{
		repo: repo,
	}
}

func (repo *ChannelRepository) tableName() (string) {
	modelValue := reflect.Indirect(reflect.ValueOf(&model.Channel{}))
	modelType := modelValue.Type()
	return repo.repo.Table(modelType)
}

func (repo *ChannelRepository) GetNewest(userID string) (*model.Channel, error) {
	rows, err := squirrel.Select("*").
		From(repo.tableName()).
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		Limit(1).
		RunWith(repo.repo.DB()).
		Query()

	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to find newest channel")
	}

	parsed := repo.repo.Parse(&model.Channel{}, rows)

	if len(parsed) != 1 {
		return nil, errors.New("could not find newest channel")
	}

	channel := parsed[0].(model.Channel)

	return &channel, nil
}

func (repo *ChannelRepository) GetBySlug(slug string) (*model.Channel, error) {
	rows, err := squirrel.Select("*").
		From(repo.tableName()).
		Where(squirrel.Eq{"slug": slug}).
		RunWith(repo.repo.DB()).
		Query()

	if err != nil {
		log.Print(err)
		return nil, err
	}

	parsed := repo.repo.Parse(&model.Channel{}, rows)
	if len(parsed) != 1 {
		return nil, nil
	}

	channel := parsed[0].(model.Channel)
	return &channel, nil
}

func (repo *ChannelRepository) GetByID(id string) (*model.Channel, error) {
	channel := &model.Channel{}

	found, err := repo.repo.FindByID(id, channel)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return channel, nil
}

func (repo *ChannelRepository) Save(channel *model.Channel) (error) {
	return repo.repo.Save(channel)
}
