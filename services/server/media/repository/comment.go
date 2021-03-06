package repository

import (
	"services/server/core/repository"
	"services/server/media/model"
	"github.com/labstack/gommon/log"
	"github.com/Masterminds/squirrel"
)

type CommentRepository struct {
	repo *repository.Repository
}

func NewCommentRepository(repo *repository.Repository) (*CommentRepository) {
	return &CommentRepository{
		repo: repo,
	}
}

func (repo *CommentRepository) GetByID(id string) (*model.Comment, error) {
	var comment model.Comment

	_, err := repo.repo.FindByID(id, &comment)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &comment, nil;
}

func (repo *CommentRepository) GetByMediaID(mediaID string) ([]model.Comment, error) {
	rows, err := squirrel.Select("*").
		From("comment").
		Where(squirrel.Eq{"media_id": mediaID}).
		OrderBy("created_at DESC").
		RunWith(repo.repo.DB()).
		Query()

	if err != nil {
		log.Print(err)
		return nil, err
	}

	// build the map of id => comment
	parsed := repo.repo.Parse(&model.Comment{}, rows)
	comments := make([]model.Comment, len(parsed))
	for i, p := range parsed {
		comments[i] = p.(model.Comment)
	}

	return comments, nil
}

func (repo *CommentRepository) Save(comment *model.Comment) (error) {
	return repo.repo.Save(comment)
}
