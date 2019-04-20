package service

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"os"
	"strings"
	"services/server/core/service/aws"
	"services/server/core/util"
)

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

func GetHomeMediaResponseProps(c *gin.Context) (gin.H) {
	infos, err := GetMediaData(util.GetSelectPage(c), nil)

	var props = gin.H{
		"error": err,
		"media": []MediaResponse{},
	}

	if err == nil {
		props["media"] = ConvertMediaData(infos, nil)
	}

	return props
	return nil
}

func GetListMediaResponseProps(c *gin.Context, data []MediaData) (gin.H) {
	uploads := []MediaResponse{}

	// there is a chance that the lambda has not started the job processing yet, so the media info won't exist
	// in those cases, we should append a pending upload in its place
	pending := getPendingUploadIfNeeded(c, data)
	if pending != nil {
		uploads = append(uploads, *pending)
	}

	uploads = append(uploads, ConvertMediaData(data, func(datum *MediaData, resp *MediaResponse) {
		resp.TranscodingStatus = aws.GetETService().GetJobStatus(datum.Media.JobID)
		resp.CanPublish = datum.Media.CanPublish()
		resp.IsPublished = datum.Media.Published
	})...)

	return gin.H{
		"uploads": uploads,
	}
}

type ConvertMediaCallback func(data *MediaData, mediaResponse *MediaResponse);

func ConvertMediaData(data []MediaData, callback ConvertMediaCallback) ([]MediaResponse) {
	media := []MediaResponse{}

	baseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_PROCESSED_BUCKET"))
	thumbBaseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_THUMBNAIL_BUCKET"))

	for _, datum := range data {
		resp := ConvertSingleMediaInfo(datum, baseUrl, thumbBaseUrl, callback)

		media = append(media, resp)
	}

	return media
}

func ConvertSingleMediaInfo(data MediaData, baseUrl string, thumbBaseUrl string, callback ConvertMediaCallback) (MediaResponse) {
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

	userId := data.Media.UserID
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
					Url:    fmt.Sprintf("%s/%s/%s/playlist.m3u8", baseUrl, userId, videoId),
				},
				{
					Type:   "mp4",
					Width:  track.Width,
					Height: track.Height,
					Url:    fmt.Sprintf("%s/%s/%s/g-720p.mp4", baseUrl, userId, videoId),
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
			fmt.Sprintf("%s/%s/%s/g-720p.mp4-00001.png", thumbBaseUrl, userId, videoId),
			fmt.Sprintf("%s/%s/%s/hls-v-1-5m-00001.png", thumbBaseUrl, userId, videoId),
			fmt.Sprintf("%s/%s/%s/hls-v-1m-00001.png", thumbBaseUrl, userId, videoId),
			fmt.Sprintf("%s/%s/%s/hls-v-1m-00001.png", thumbBaseUrl, userId, videoId),
			fmt.Sprintf("%s/%s/%s/hls-v-400k-00001.png", thumbBaseUrl, userId, videoId),
			fmt.Sprintf("%s/%s/%s/hls-v-600k-00001.png", thumbBaseUrl, userId, videoId),
		},
		Videos: videos,
	}

	// optionally allow additional data to be added to the info
	if callback != nil {
		callback(&data, &resp)
	}

	return resp;
}

func getPendingUploadIfNeeded(c *gin.Context, data []MediaData) (*MediaResponse) {
	params := c.Request.URL.Query()
	bucket, okBucket := params["bucket"]
	key, okKey := params["key"]
	if okBucket && okKey {
		pieces := strings.Split(key[0], "/")
		videoId := pieces[len(pieces)-1]
		canAddPending := len(data) == 0 || data[0].Media.ID != videoId

		// basically, we don't have the latest info from the trannscoder, but the file was definitely uploaded
		// so we should append the info anyways...
		if canAddPending && aws.GetS3Service().FileExists(bucket[0], key[0]) {
			return &MediaResponse{
				ID:                videoId,
				CanPublish:        false,
				Title:             "",
				Description:       "",
				TranscodingStatus: "pending",
				Hashtags:          []string{},
				IsPublished:       false,
				Thumbnails:        []string{},
				Videos:            []MediaVideoResponse{},
			}
		}
	}
	return nil;
}
