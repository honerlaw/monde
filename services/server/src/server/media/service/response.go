package service

import (
	"github.com/gin-gonic/gin"
	"server/core/util"
	"fmt"
	"os"
)

type MediaVideoResponse struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type MediaResponse struct {
	ID          string               `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Hashtags    []string             `json:"hashtags"`
	Thumbnails  []string             `json:"thumbnails"`
	Videos      []MediaVideoResponse `json:"videos"`
}

// @todo we should do something very similar for /media/list
func GetHomeMediaResponseProps(c *gin.Context) (gin.H) {
	infos, err := GetMediaInfo(util.GetSelectPage(c))

	var props = gin.H{
		"error": err,
		"media": []MediaResponse{},
	}

	if err != nil {

		baseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_PROCESSED_BUCKET"))
		thumbBaseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s/", os.Getenv("AWS_THUMBNAIL_BUCKET"))

		for _, info := range *infos {

			hashtags := make([]string, 0)
			for _, hashtag := range info.Hashtags {
				hashtags = append(hashtags, hashtag.Tag)
			}

			props["media"] = append(props["media"].([]MediaResponse), MediaResponse{
				ID: info.VideoID,
				Title: info.Title,
				Description: info.Description,
				Hashtags: hashtags,
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
			})
		}

	}

	return props
}
