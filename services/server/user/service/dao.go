package service

import (
	"services/server/core/repository"
	"log"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"services/server/user/model"
	"strings"
	"time"
)

func FindUserByUsername(username string) (*model.User) {
	rows, err := squirrel.
		Select(strings.Join(model.UserColumns, ",")).
		From("user").
		Where(squirrel.Eq{"username": username}).
		RunWith(repository.GetRepository().DB()).
		Query()

	if err != nil {
		log.Print(err)
		return nil
	}

	users := parse(rows)

	if len(users) == 0 {
		return nil
	}

	return &users[0]
}

func FindUserByID(id uint) (*model.User) {
	rows, err := squirrel.
		Select(strings.Join(model.UserColumns, ",")).
		From("user").
		Where(squirrel.Eq{"id": id}).
		RunWith(repository.GetRepository().DB()).
		Query()

	if err != nil {
		log.Print(err)
		return nil
	}

	users := parse(rows)

	if len(users) == 0 {
		return nil
	}

	return &users[0]
}

func SaveUser(user *model.User) (error) {
	found := FindUserByID(user.ID)
	if found == nil {

		user.CreatedAt = time.Now()
		user.UpdatedAt = user.CreatedAt

		_, err := squirrel.Insert("user").
			Columns("created_at", "updated_at", "username", "hash").
			Values(user.CreatedAt, user.UpdatedAt, user.Username, user.Hash).
			RunWith(repository.GetRepository().DB()).
			Query()

		if err != nil {
			log.Print(err)
			return err
		}

		return nil
	}

	_, err := squirrel.Update("user").
		Set("updated_at", time.Now()).
		RunWith(repository.GetRepository().DB()).
		Query()

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func parse(rows *sql.Rows) ([]model.User) {
	defer rows.Close()

	users := []model.User{}
	rows.Columns()

	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.Username, &user.Hash)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return users;
}