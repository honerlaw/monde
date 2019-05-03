package repository

import (
	"services/server/core/repository"
	"services/server/user/model"
)

type UserRepository struct {
	repo *repository.Repository
}

func NewUserRepository(repo *repository.Repository) (*UserRepository) {
	return &UserRepository{
		repo: repo,
	}
}

func (repo *UserRepository) FindByID(id string) (*model.User) {
	user := &model.User{}
	found, err := repo.repo.FindByID(id, user)
	if !found || err != nil {
		return nil
	}
	return user;
}

func (repo *UserRepository) Save(user *model.User) (error) {
	return repo.repo.Save(user)
}
