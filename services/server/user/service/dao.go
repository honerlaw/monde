package service

import (
	"services/server/core/repository"
	"log"
	"github.com/Masterminds/squirrel"
	"services/server/user/model"
)

func FindUserByUsername(username string) (*model.User) {
	rows, err := squirrel.
		Select("*").
		From("user").
		Where(squirrel.Eq{"username": username}).
		RunWith(repository.GetRepository().DB()).
		Query()

	if err != nil {
		log.Print(err)
		return nil
	}

	users := repository.GetRepository().Parse(&model.User{}, rows)
	if len(users) == 0 {
		return nil
	}
	user := users[0].(model.User)

	return &user
}
