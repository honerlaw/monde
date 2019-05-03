package service

import (
	"golang.org/x/crypto/bcrypt"
	"errors"
	"services/server/user/model"
	"services/server/user/repository"
	repository2 "services/server/core/repository"
	"github.com/labstack/gommon/log"
	"strings"
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
	contactService *ContactService
	channelService *ChannelService
	userRepository *repository.UserRepository
}

func NewUserService(contactService *ContactService, channelService *ChannelService, userRepository *repository.UserRepository) (*UserService) {
	return &UserService{
		contactService: contactService,
		channelService: channelService,
		userRepository: userRepository,
	}
}

func (service *UserService) FindUserByEmail(email string) (*model.User) {
	contact, err := service.contactService.FindByEmail(email)

	if contact == nil || err != nil {
		return nil
	}

	user := service.userRepository.FindByID(contact.UserID)
	if user == nil {
		return nil
	}

	return user
}

func (service *UserService) Verify(req VerifyRequest) (*model.User, error) {
	user := service.FindUserByEmail(req.Email)
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
	if !service.contactService.IsPotentiallyValidContact(req.Email, "email") {
		return nil, errors.New("invalid email address")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters in length")
	}

	if req.Password != req.VerifyPassword {
		return nil, errors.New("passwords do not match")
	}

	user := service.FindUserByEmail(req.Email)
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

	user = &model.User{
		Hash:     string(hash),
	}

	err = service.userRepository.Save(user)
	if err != nil {
		tx.Rollback()
		return nil, genericError
	}

	// the default channel will just be the portion of the email before the @ sign
	title := strings.Split(req.Email, "@")[0]
	_, err = service.channelService.Create(user.ID, title)
	if err != nil {
		tx.Rollback()
		return nil, genericError
	}

	_, err = service.contactService.Create(user.ID, req.Email, "email")
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
