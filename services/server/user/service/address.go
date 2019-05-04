package service

import (
	"services/server/user/repository"
	"github.com/pkg/errors"
	"services/server/user/model"
)

type AddressData struct {
	addresses []model.Address
}

type AddressCreateRequest struct {
	Type    string `form:"type" json:"type" binding:"required"`
	LineOne string `form:"line_one" json:"line_one" binding:"required"`
	LineTwo string `form:"line_two" json:"line_two"`
	City    string `form:"city" json:"city" binding:"required"`
	State   string `form:"state" json:"state" binding:"required"`
	ZipCode string `form:"zip_code" json:"zip_code" binding:"required"`
	Country string `form:"country" json:"country" binding:"required"`
}

type AddressService struct {
	addressRepository *repository.AddressRepository
}

func NewAddressService(addressRepository *repository.AddressRepository) (*AddressService) {
	return &AddressService{
		addressRepository: addressRepository,
	}
}

func (service *AddressService) GetAddressesByUserID(userID string) ([]model.Address) {
	addresses, err := service.addressRepository.FindByUserID(userID)
	if err != nil {
		return nil
	}

	return addresses
}

func (service *AddressService) Create(userID string, req *AddressCreateRequest) (*model.Address, error) {
	if req.Type != "home" && req.Type != "business" {
		return nil, errors.New("invalid address type")
	}

	// @todo we should verify these fields in some capacity
	address := &model.Address{
		UserID:  userID,
		Type:    req.Type,
		LineOne: req.LineOne,
		LineTwo: req.LineTwo,
		City:    req.City,
		State:   req.State,
		ZipCode: req.ZipCode,
		Country: req.Country,
	}

	err := service.addressRepository.Save(address)
	if err != nil {
		return nil, errors.New("something went wrong")
	}
	return address, nil
}
