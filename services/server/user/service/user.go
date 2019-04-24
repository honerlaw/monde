package service

import (
	"golang.org/x/crypto/bcrypt"
	"errors"
	"services/server/user/model"
	"services/server/user/repository"
	repository2 "services/server/core/repository"
	"github.com/labstack/gommon/log"
)

type VerifyRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type CreateRequest struct {
	VerifyRequest
	VerifyPassword string `form:"verify_password" binding:"required"`
}

type UserService struct {
	channelService *ChannelService
	userRepository *repository.UserRepository
}

func NewUserService(channelService *ChannelService, userRepository *repository.UserRepository) (*UserService) {
	return &UserService{
		channelService: channelService,
		userRepository: userRepository,
	}
}

func (service *UserService) Verify(req VerifyRequest) (*model.User, error) {
	user := service.userRepository.FindByUsername(req.Username)

	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (service *UserService) Create(req CreateRequest) (*model.User, error) {
	if len(req.Username) < 6 {
		return nil, errors.New("username must be at least 6 characters in length")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters in length")
	}

	if req.Password != req.VerifyPassword {
		return nil, errors.New("passwords do not match")
	}

	user := service.userRepository.FindByUsername(req.Username)
	if user != nil {
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return nil, errors.New("something went wrong. please try again")
	}

	genericError := errors.New("something went wrong. please try again")

	tx, err := repository2.GetRepository().DB().Begin()
	if err != nil {
		log.Print(err)
		return nil, genericError
	}

	user = &model.User{
		Username: req.Username,
		Hash:     string(hash),
	}

	err = service.userRepository.Save(user)
	if err != nil {
		tx.Rollback()
		return nil, genericError
	}

	_, err = service.channelService.Create(user.ID, user.Username)
	if err != nil {
		tx.Rollback()
		return nil, genericError
	}

	err = tx.Commit()
	if err != nil {
		return nil, genericError
	}

	return user, nil
}
