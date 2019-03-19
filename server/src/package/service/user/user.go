package user

import (
	"package/model"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"package/repository"
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
	var user model.User
	repository.DB.Where(model.User{Username: req.Username}).First(&user)

	// can't check an empty struct, so just make sure their username exists to see if we found it
	if user.Username == "" {
		return nil, errors.New("invalid username or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
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

	var user model.User
	repository.DB.Where(model.User{Username: req.Username}).First(&user)

	// check the username since we can't easily compare to an empty struct
	if user.Username != "" {
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
