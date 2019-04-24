package service

import (
	"services/server/user/repository"
	"services/server/user/model"
)

type ChannelService struct {
	channelRepository *repository.ChannelRepository
}

func NewChannelService(channelRepository *repository.ChannelRepository) (*ChannelService) {
	return &ChannelService{
		channelRepository: channelRepository,
	}
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
	slug := title
	channel := model.Channel{
		UserID: userID,
		Title:  title,
		Slug:   slug,
	}

	err := service.channelRepository.Save(&channel)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

