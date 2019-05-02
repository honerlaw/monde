package repository

import (
	"services/server/core/repository"
	"services/server/user/model"
	"github.com/Masterminds/squirrel"
	"log"
)

type UserRepository struct {
	repo *repository.Repository
}

func NewUserRepository(repo *repository.Repository) (*UserRepository) {
	return &UserRepository{
		repo: repo,
	}
}

func (repo *UserRepository) FindByEmail(email string) (*model.User) {
	rows, err := squirrel.
		Select("*").
		From("user").
		Where(squirrel.Eq{"email": email}).
		RunWith(repository.GetRepository().DB()).
		Query()

	if err != nil {
		log.Print(err)
		return nil
	}

	var user model.User
	users := repository.GetRepository().Parse(&user, rows)
	if len(users) == 0 {
		return nil
	}
	user = users[0].(model.User)

	return &user
}

func (repo *UserRepository) Save(user *model.User) (error) {
	return repo.repo.Save(user)
}
