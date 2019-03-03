package service

import (
	"package/model"
	"package/repository"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

type VerifyUserRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type CreateUserRequest struct {
	VerifyUserRequest
	VerifyPassword string `form:"verify_password" binding:"required"`
}

func VerifyUser(req VerifyUserRequest) (*model.User, error) {
	var user model.User
	repository.DB.Where(model.User{Username: req.Username}).First(&user)

	if (model.User{}) == user {
		return nil, errors.New("invalid username or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}

func CreateUser(req CreateUserRequest) (*model.User, error) {
	if len(req.Username) < 6 {
		return nil, errors.New("username must be at least 6 characters in length")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters in length")
	}

	if req.Password != req.VerifyPassword {
		return nil, errors.New("passwords do not match")
	}

	var user model.User
	repository.DB.Where(model.User{Username: req.Username}).First(&user)

	if (model.User{}) != user {
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return nil, errors.New("something went wrong. please try again")
	}

	user = model.User{
		Username: req.Username,
		Hash:     string(hash),
	}

	repository.DB.Save(&user);
	return &user, nil
}
