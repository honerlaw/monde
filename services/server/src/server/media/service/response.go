package service

import (
	"github.com/gin-gonic/gin"
	"server/core/util"
	"fmt"
	"os"
	"strings"
	"server/core/service/aws"
	"server/media/model"
)

type MediaVideoResponse struct {
	Type string `json:"type"`
	Url  string `json:"url"`
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
	infos, err := GetMediaInfo(util.GetSelectPage(c))

	var props = gin.H{
		"error": err,
		"media": []MediaResponse{},
	}

	if err != nil {
		props["media"] = ConvertMediaInfo(infos, nil)
	}

	return props
}

func GetListMediaResponseProps(c *gin.Context, infos *[]model.MediaInfo) (gin.H) {
	uploads := []MediaResponse{}

	// there is a chance that the lambda has not started the job processing yet, so the media info won't exist
	// in those cases, we should append a pending upload in its place
	pending := getPendingUploadIfNeeded(c, infos)
	if pending != nil {
		uploads = append(uploads, *pending)
	}

	uploads = append(uploads, *ConvertMediaInfo(infos, func(info *model.MediaInfo, resp *MediaResponse) {
		resp.TranscodingStatus = aws.GetETService().GetJobStatus(info.JobID)
		resp.CanPublish = info.CanPublish()
		resp.IsPublished = info.Published
	})...)

	return gin.H{
		"uploads": uploads,
	}
}

func ConvertMediaInfo(infos *[]model.MediaInfo, callback func(info *model.MediaInfo, mediaResponse *MediaResponse)) (*[]MediaResponse) {
	media := []MediaResponse{}

	baseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_PROCESSED_BUCKET"))
	thumbBaseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s/", os.Getenv("AWS_THUMBNAIL_BUCKET"))

	for _, info := range *infos {

		hashtags := make([]string, 0)
		for _, hashtag := range info.Hashtags {
			hashtags = append(hashtags, hashtag.Tag)
		}

		resp := MediaResponse{
			ID:          info.VideoID,
			Title:       info.Title,
			Description: info.Description,
			Hashtags:    hashtags,
			Thumbnails: []string{
				fmt.Sprintf("%s/%s/%s/g-720p.mp4-00001.png", thumbBaseUrl, info.UserID, info.VideoID),
				fmt.Sprintf("%s/%s/%s/hls-v-1-5m-00001.png", thumbBaseUrl, info.UserID, info.VideoID),
				fmt.Sprintf("%s/%s/%s/hls-v-1m-00001.png", thumbBaseUrl, info.UserID, info.VideoID),
				fmt.Sprintf("%s/%s/%s/hls-v-1m-00001.png", thumbBaseUrl, info.UserID, info.VideoID),
				fmt.Sprintf("%s/%s/%s/hls-v-400k-00001.png", thumbBaseUrl, info.UserID, info.VideoID),
				fmt.Sprintf("%s/%s/%s/hls-v-600k-00001.png", thumbBaseUrl, info.UserID, info.VideoID),
			},
			Videos: []MediaVideoResponse{
				{
					Type: "hls",
					Url:  fmt.Sprintf("%s/%s/%s/playlist.m3u8", baseUrl, info.UserID, info.VideoID),
				},
				{
					Type: "mp4",
					Url:  fmt.Sprintf("%s/%s/%s/g-720p.mp4", baseUrl, info.UserID, info.VideoID),
				},
			},
		}

		// optionally allow additional data to be added to the info
		if callback != nil {
			callback(&info, &resp)
		}

		media = append(media, resp)
	}

	return &media
}

func getPendingUploadIfNeeded(c *gin.Context, infos *[]model.MediaInfo) (*MediaResponse) {
	params := c.Request.URL.Query()
	bucket, okBucket := params["bucket"]
	key, okKey := params["key"]
	if okBucket && okKey {
		pieces := strings.Split(key[0], "/")
		videoId := pieces[len(pieces)-1]
		info := (*infos)[0]

		// basically, we don't have the latest info from the trannscoder, but the file was definitely uploaded
		// so we should append the info anyways...
		if info.VideoID != videoId && aws.GetS3Service().FileExists(bucket[0], key[0]) {
			return &MediaResponse{
				ID:                videoId,
				CanPublish:        false,
				Title:             nil,
				Description:       nil,
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
