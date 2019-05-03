package service

import (
	"services/server/user/repository"
	"services/server/user/model"
	"github.com/pkg/errors"
	"strings"
	"regexp"
	"services/server/core/service/aws"
	"math/rand"
)

type ContactService struct {
	contactRepository *repository.ContactRepository
	emailRegex        *regexp.Regexp
	phoneRegex        *regexp.Regexp
}

func NewContactService(contactRepository *repository.ContactRepository) (*ContactService) {
	return &ContactService{
		contactRepository: contactRepository,
		phoneRegex:        regexp.MustCompile("[^0-9]+"),
		emailRegex:        regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"),
	}
}

func (service *ContactService) NormalizeContact(contact string, contactType string) (string) {
	if contactType == "sms" {
		// basically only leave the numbers
		return service.phoneRegex.ReplaceAllString(contact, "")
	}
	return strings.TrimSpace(strings.ToLower(contact))
}

func (service *ContactService) IsPotentiallyValidContact(contact string, contactType string) (bool) {
	contact = service.NormalizeContact(contact, contactType)
	if contactType == "sms" {
		// this includes country codes
		return len(contact) < 8 || len(contact) > 16
	}
	return service.emailRegex.MatchString(contact)
}

func (service *ContactService) FindByEmail(contact string) (*model.Contact, error) {
	contact = service.NormalizeContact(contact, "email")
	return service.contactRepository.FindByContact(contact)
}

func (service *ContactService) Create(userID string, contact string, contactType string) (*model.Contact, error) {
	if contactType != "email" && contactType != "sms" {
		return nil, errors.New("invalid contact method")
	}

	if !service.IsPotentiallyValidContact(contact, contactType) {
		return nil, errors.New("invalid contact")
	}

	contact = service.NormalizeContact(contact, contactType)

	found, err := service.contactRepository.FindByContact(contact)
	if err != nil || found != nil {
		return nil, errors.New("contact already exists")
	}

	contactModel := &model.Contact{
		UserID: userID,
		Contact: contact,
		Type: contactType,
		Code: service.generateCode(),
		Verified: false,
	}

	err = service.contactRepository.Save(contactModel)
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	// send email notificationn
	// @todo send out a verification link
	err = aws.GetSESService().SendEmail(contactModel.Contact, "Testing", "Testing!")
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	return contactModel, nil
}

func (service *ContactService) generateCode() string {
	length := 6
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(65 + rand.Intn(25))  //A=65 and Z = 65+25
	}
	return string(bytes)
}