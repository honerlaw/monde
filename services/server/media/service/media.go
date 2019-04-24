package service

import (
	"time"
	"strings"
	"regexp"
	"services/server/media/model"
	"errors"
	"services/server/media/repository"
	"fmt"
	"os"
	"services/server/core/util"
	"github.com/Masterminds/squirrel"
)

type UpdateRequest struct {
	VideoID     string `form:"video_id" binding:"required"`
	Title       string `form:"title"`
	Description string `form:"description"`
	Hashtags    string `form:"hashtags"`
}

type PublishRequest struct {
	VideoID string `form:"video_id" binding:"required"`
}

type MediaVideoResponse struct {
	Type   string `json:"type"`
	Url    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

type MediaResponse struct {
	ID                string               `json:"id"`
	Title             string               `json:"title"`
	Description       string               `json:"description"`
	Hashtags          []string             `json:"hashtags"`
	Thumbnails        []string             `json:"thumbnails"`
	Videos            []MediaVideoResponse `json:"videos"`
	TranscodingStatus string               `json:"transcoding_status"`
	CanPublish        bool                 `json:"can_publish"`
	IsPublished       bool                 `json:"is_published"`
}

type ConvertMediaCallback func(data *repository.MediaData, mediaResponse *MediaResponse);

type MediaService struct {
	mediaRepository *repository.MediaRepository
	hashtagRepository *repository.HashtagRepository
}

func NewMediaService(mediaRepository *repository.MediaRepository, hashtagRepository *repository.HashtagRepository) (*MediaService) {
	return &MediaService{
		mediaRepository: mediaRepository,
		hashtagRepository: hashtagRepository,
	}
}

func (service *MediaService) List(selectPage *util.SelectPage) ([]repository.MediaData, error) {
	return service.mediaRepository.List(selectPage, squirrel.Eq{
		"published": true,
	})
}

func (service *MediaService) GetByChannelID(channelID string, selectPage *util.SelectPage) ([]repository.MediaData, error) {
	return service.mediaRepository.GetByChannelID(channelID, selectPage)
}

func (service *MediaService) GetByVideoID(videoID string) (*repository.MediaData, error) {
	return service.mediaRepository.GetByVideoID(videoID)
}

func (service *MediaService) Update(req UpdateRequest) (error) {
	data, err := service.mediaRepository.GetByVideoID(req.VideoID)
	if err != nil {
		return err
	}

	data.Media.Title = req.Title
	data.Media.Description = req.Description

	// if they remove the description, unpublish the video
	if len(strings.TrimSpace(data.Media.Description)) == 0 {
		data.Media.Published = false;
	}

	var hashtags []model.Hashtag
	regex, _ := regexp.Compile("^#\\w+")
	tags := strings.Split(req.Hashtags, " ")
	for _, tag := range tags {
		trimmedTag := strings.TrimSpace(tag)
		if regex.MatchString(trimmedTag) && len(trimmedTag) < 75 && len(trimmedTag) > 1 {
			// fetch the existing hashtag if it exists
			hashtag, err := service.hashtagRepository.Get(trimmedTag)
			if hashtag == nil || err != nil {
				hashtag = &model.Hashtag{
					Tag: trimmedTag,
				}
			}

			hashtags = append(hashtags, *hashtag)
		}
	}
	data.Tags = hashtags

	err = service.mediaRepository.Save(data)
	if err != nil {
		return err
	}

	return nil
}

func (service *MediaService) TogglePublish(req PublishRequest) (error) {
	data, err := service.mediaRepository.GetByVideoID(req.VideoID)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(data.Media.Description)) == 0 {
		return errors.New("a description is required to publish videos")
	}

	data.Media.Published = !data.Media.Published
	if data.Media.Published {
		data.Media.PublishedDate = time.Now()
	}

	err = service.mediaRepository.Save(data)
	if err != nil {
		return err
	}

	return nil
}

func (service *MediaService) ConvertMediaData(data []repository.MediaData, callback ConvertMediaCallback) ([]MediaResponse) {
	media := []MediaResponse{}

	baseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_PROCESSED_BUCKET"))
	thumbBaseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_THUMBNAIL_BUCKET"))

	for _, datum := range data {
		resp := service.ConvertSingleMediaInfo(datum, baseUrl, thumbBaseUrl, callback)

		media = append(media, resp)
	}

	return media
}

func (service *MediaService) ConvertSingleMediaInfo(data repository.MediaData, baseUrl string, thumbBaseUrl string, callback ConvertMediaCallback) (MediaResponse) {
	if baseUrl == "" {
		baseUrl = fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_PROCESSED_BUCKET"))
	}
	if thumbBaseUrl == "" {
		thumbBaseUrl = fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_THUMBNAIL_BUCKET"))
	}

	hashtags := make([]string, 0)
	for _, hashtag := range data.Tags {
		hashtags = append(hashtags, hashtag.Tag)
	}

	channelId := data.Media.ChannelID
	videoId := data.Media.ID

	videos := []MediaVideoResponse{}
	for _, track := range data.Tracks {
		if track.Type == "Video" {
			// @todo
			// we need post processing information about the videos (e.g. we need to store the types of videos
			// associated files, genenral video information, etc
			videos = append(videos, []MediaVideoResponse{
				{
					Type:   "hls",
					Width:  track.Width,
					Height: track.Height,
					Url:    fmt.Sprintf("%s/%s/%s/playlist.m3u8", baseUrl, channelId, videoId),
				},
				{
					Type:   "mp4",
					Width:  track.Width,
					Height: track.Height,
					Url:    fmt.Sprintf("%s/%s/%s/g-720p.mp4", baseUrl, channelId, videoId),
				},
			}...)
		}
	}

	resp := MediaResponse{
		ID:          data.Media.ID,
		Title:       data.Media.Title,
		Description: data.Media.Description,
		Hashtags:    hashtags,
		Thumbnails: []string{
			fmt.Sprintf("%s/%s/%s/g-720p.mp4-00001.png", thumbBaseUrl, channelId, videoId),
			fmt.Sprintf("%s/%s/%s/hls-v-1-5m-00001.png", thumbBaseUrl, channelId, videoId),
			fmt.Sprintf("%s/%s/%s/hls-v-1m-00001.png", thumbBaseUrl, channelId, videoId),
			fmt.Sprintf("%s/%s/%s/hls-v-1m-00001.png", thumbBaseUrl, channelId, videoId),
			fmt.Sprintf("%s/%s/%s/hls-v-400k-00001.png", thumbBaseUrl, channelId, videoId),
			fmt.Sprintf("%s/%s/%s/hls-v-600k-00001.png", thumbBaseUrl, channelId, videoId),
		},
		Videos: videos,
	}

	// optionally allow additional data to be added to the info
	if callback != nil {
		callback(&data, &resp)
	}

	return resp;
}
