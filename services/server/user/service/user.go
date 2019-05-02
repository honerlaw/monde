package service

import (
	"golang.org/x/crypto/bcrypt"
	"errors"
	"services/server/user/model"
	"services/server/user/repository"
	repository2 "services/server/core/repository"
	"github.com/labstack/gommon/log"
	"regexp"
	"strings"
	"services/server/core/service/aws"
)

type VerifyRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type CreateRequest struct {
	VerifyRequest
	VerifyPassword string `form:"verify_password" binding:"required"`
}

type UserService struct {
	channelService *ChannelService
	userRepository *repository.UserRepository
	emailRegex     *regexp.Regexp
}

func NewUserService(channelService *ChannelService, userRepository *repository.UserRepository) (*UserService) {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return &UserService{
		channelService: channelService,
		userRepository: userRepository,
		emailRegex: emailRegex,
	}
}

func (service *UserService) Verify(req VerifyRequest) (*model.User, error) {
	user := service.userRepository.FindByEmail(req.Email)

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (service *UserService) Create(req CreateRequest) (*model.User, error) {
	if !service.emailRegex.MatchString(req.Email) {
		return nil, errors.New("invalid email address")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters in length")
	}

	if req.Password != req.VerifyPassword {
		return nil, errors.New("passwords do not match")
	}

	user := service.userRepository.FindByEmail(req.Email)
	if user != nil {
		return nil, errors.New("user already exists")
	}

	// generic error to use for several scenarios
	genericError := errors.New("something went wrong. please try again")

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return nil, genericError
	}

	tx, err := repository2.GetRepository().DB().Begin()
	if err != nil {
		log.Print(err)
		return nil, genericError
	}

	// usernname is the first portion before the @sign
	username := strings.Split(req.Email, "@")[0]

	user = &model.User{
		Email: req.Email,
		Username: username,
		Hash:  string(hash),
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

	err = aws.GetSESService().SendEmail(req.Email, "Testing", "Testing!")
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
