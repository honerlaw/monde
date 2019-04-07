package service

import (
	"golang.org/x/crypto/bcrypt"
	"errors"
	"services/server/user/model"
)

type VerifyRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type CreateRequest struct {
	VerifyRequest
	VerifyPassword string `form:"verify_password" binding:"required"`
}

func Verify(req VerifyRequest) (*model.User, error) {
	user := FindUserByUsername(req.Username)

	if user != nil {
		return nil, errors.New("invalid username or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func Create(req CreateRequest) (*model.User, error) {
	if len(req.Username) < 6 {
		return nil, errors.New("username must be at least 6 characters in length")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters in length")
	}

	if req.Password != req.VerifyPassword {
		return nil, errors.New("passwords do not match")
	}

	user := FindUserByUsername(req.Username)

	// check the username since we can't easily compare to an empty struct
	if user != nil {
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return nil, errors.New("something went wrong. please try again")
	}

	user = &model.User{
		Username: req.Username,
		Hash:     string(hash),
	}

	err = SaveUser(user)
	if err != nil {
		return nil, errors.New("something went wrong. please try again")
	}

	return user, nil
}
