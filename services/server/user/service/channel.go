package service

import (
	"services/server/user/repository"
	"services/server/user/model"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

type ChannelService struct {
	channelRepository *repository.ChannelRepository
}

func NewChannelService(channelRepository *repository.ChannelRepository) (*ChannelService) {
	return &ChannelService{
		channelRepository: channelRepository,
	}
}

func (service *ChannelService) GetBySlug(slug string) (*model.Channel, error) {
	return service.channelRepository.GetBySlug(slug)
}

func (service *ChannelService) GetByUserID(userID string, id *string) (*model.Channel, error) {
	if id == nil {
		return service.channelRepository.GetNewest(userID)
	}
	return service.GetByID(*id)
}

func (service *ChannelService) GetByID(channelID string) (*model.Channel, error) {
	return service.channelRepository.GetByID(channelID)
}

func (service *ChannelService) Create(userID string, title string) (*model.Channel, error) {
	s := slug.Make(title)

	// @todo we should attempt to generate a slug multiple times if we need to
	found, err := service.GetBySlug(s)
	if err != nil {
		return nil, err
	}
	if found != nil {
		return nil, errors.New("slug already exists for given title")
	}


	channel := model.Channel{
		UserID: userID,
		Title:  title,
		Slug:   s,
	}

	err = service.channelRepository.Save(&channel)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

