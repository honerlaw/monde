package repository

import (
	"services/server/core/repository"
	"services/server/media/model"
	"github.com/Masterminds/squirrel"
	"errors"
)

type HashtagRepository struct {
	repo *repository.Repository
}

func NewHashtagRepository(repo *repository.Repository) (*HashtagRepository) {
	return &HashtagRepository{
		repo: repo,
	}
}

func (repo *HashtagRepository) Get(tag string) (*model.Hashtag, error) {
	rows, err := squirrel.Select("*").
		From("hashtag").
		Where(squirrel.Eq{"tag": tag}).
		RunWith(repo.repo.DB()).Query()

	if err != nil {
		return nil, errors.New("failed to find hashtag")
	}

	parsed := repo.repo.Parse(&model.Hashtag{}, rows)
	models := make([]model.Hashtag, len(parsed))
	for i, p := range parsed {
		models[i] = p.(model.Hashtag)
	}

	if len(models) == 0 {
		return nil, errors.New("failed to find hashtag")
	}

	return &models[0], nil
}