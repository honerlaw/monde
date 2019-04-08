package service

import (
	"services/server/core/repository"
	"log"
	"github.com/Masterminds/squirrel"
	"services/server/user/model"
	"github.com/gin-gonic/gin/json"
	"fmt"
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

	var user model.User
	users := repository.GetRepository().Parse(&user, rows)
	if len(users) == 0 {
		return nil
	}
	user = users[0].(model.User)

	data, _ := json.Marshal(user)

	fmt.Print(string(data))

	return &user
}
