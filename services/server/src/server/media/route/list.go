package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"fmt"
	"strconv"
	"server/media/service"
	"strings"
	"server/core/service/aws"
	"server/core/repository"
	"server/user/middleware"
	"server/core/render"
	"server/media/model"
	"server/core/util"
)

func List(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY").(*middleware.AuthPayload)
	props := gin.H{
		"uploads":     []gin.H{},
	}

	if payload == nil {
		render.RenderPage(c, http.StatusOK, props)
		return
	}

	// fetch requested media info for given page
	infos, err := service.GetMediaInfoByUserId(payload.ID, util.GetSelectPage(c))
	if err != nil {
		render.RenderPage(c, http.StatusInternalServerError, nil)
		return
	}

	// there is a chance that the lambda has not started the job processing yet, so the media info won't exist
	// in those cases, we should append a pending upload in its place
	pending := getPendingUploadIfNeeded(c, infos)
	if pending != nil {
		props["uploads"] = append(props["uploads"].([]gin.H), pending)
	}

	// convert all of the infos to some data needed to properly render the list of items
	for _, info := range *infos {
		status := aws.GetETService().GetJobStatus(info.JobID)

		// get all of the hash tags and convert to a string array
		hashtags := make([]string, 0) // allocate an empty array so json converts it properly
		repository.DB.Model(info).Related(&info.Hashtags, "Hashtags")
		for _, hashtag := range info.Hashtags {
			hashtags = append(hashtags, hashtag.Tag)
		}

		baseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_PROCESSED_BUCKET"))
		thumbBaseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s/", os.Getenv("AWS_THUMBNAIL_BUCKET"))

		userId := strconv.FormatUint(uint64(payload.ID), 10)

		props["uploads"] = append(props["uploads"].([]gin.H), gin.H{
			"videoId":    info.VideoID,
			"canPublish": len(strings.TrimSpace(info.Description)) > 0,
			"info": gin.H{
				"title":       info.Title,
				"description": info.Description,
				"status":      status,
				"hashtags":    hashtags,
				"published":   info.Published,
			},
			"thumbs": []string{
				fmt.Sprintf("%s/%s/%s/g-720p.mp4-00001.png", thumbBaseUrl, userId, info.VideoID),
				fmt.Sprintf("%s/%s/%s/hls-v-1-5m-00001.png", thumbBaseUrl, userId, info.VideoID),
				fmt.Sprintf("%s/%s/%s/hls-v-1m-00001.png", thumbBaseUrl, userId, info.VideoID),
				fmt.Sprintf("%s/%s/%s/hls-v-1m-00001.png", thumbBaseUrl, userId, info.VideoID),
				fmt.Sprintf("%s/%s/%s/hls-v-400k-00001.png", thumbBaseUrl, userId, info.VideoID),
				fmt.Sprintf("%s/%s/%s/hls-v-600k-00001.png", thumbBaseUrl, userId, info.VideoID),
			},
			"videos": []gin.H{
				{
					"type": "hls",
					"url":  fmt.Sprintf("%s/%s/%s/playlist.m3u8", baseUrl, userId, info.VideoID),
				},
				{
					"type": "mp4",
					"url":  fmt.Sprintf("%s/%s/%s/g-720p.mp4", baseUrl, userId, info.VideoID),
				},
			},
		})
	}

	render.RenderPage(c, http.StatusOK, props)
}

func getPendingUploadIfNeeded(c *gin.Context, infos *[]model.MediaInfo) (gin.H) {
	params := c.Request.URL.Query()
	bucket, okBucket := params["bucket"]
	key, okKey := params["key"]
	if okBucket && okKey {
		pieces := strings.Split(key[0], "/")
		videoId := pieces[len(pieces) - 1]
		info := (*infos)[0]

		// basically, we don't have the latest info from the trannscoder, but the file was definitely uploaded
		// so we should append the info anyways...
		if info.VideoID != videoId && aws.GetS3Service().FileExists(bucket[0], key[0]) {
			return gin.H{
				"videoId": videoId,
				"canPublish": false,
				"info": gin.H{
					"title": nil,
					"description": nil,
					"status": "pending",
					"hashtags": []string{},
					"published": false,
				},
				"thumbs": []gin.H{},
				"videos": []gin.H{},
			}
		}
	}
	return nil;
}