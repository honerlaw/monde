package route

import (
	"github.com/gin-gonic/gin"
	"server/util"
	"net/http"
	"server/middleware/auth"
	"server/service/aws"
	"os"
	"fmt"
	"strconv"
	"server/media/service"
	"server/repository"
)

func List(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY").(*auth.AuthPayload)
	props := gin.H{
		"authPayload": payload,
		"uploads":     []gin.H{},
	}

	if payload == nil {
		util.RenderPage(c, http.StatusOK, "UploadListPage", props)
		return
	}

	// @todo get page info from query parameters
	infos, err := service.GetMediaInfoByUserId(payload.ID, &service.SelectPage{
		Page:  0,
		Count: 50,
	})
	if err != nil {
		util.RenderPage500(c)
		return
	}

	for _, info := range *infos {
		status := aws.GetETService().GetJobStatus(info.JobID)

		// get all of the hash tags and convert to a string array
		hashtags := make([]string, 0) // allocate an empty array so json converts it properly
		repository.DB.Model(info).Related(&info.Hashtags, "Hashtags")
		for _, hashtag := range info.Hashtags {
			hashtags = append(hashtags, hashtag.Tag)
		}

		baseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_PROCESSED_BUCKET"))
		// thumbBaseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s/", os.Getenv("AWS_THUMBNAIL_BUCKET"))

		props["uploads"] = append(props["uploads"].([]gin.H), gin.H{
			"videoId": info.VideoID,
			"info": gin.H{
				"title":       info.Title,
				"description": info.Description,
				"status":      status,
				"hashtags":    hashtags,
				"published":   info.Published,
			},
			"thumbs": []string{""}, // @todo figure out what this actually is
			"videos": []gin.H{
				{
					"type": "hls",
					"url":  fmt.Sprintf("%s/%s/%s/playlist.m3u8", baseUrl, strconv.FormatUint(uint64(payload.ID), 10), info.VideoID),
				},
				{
					"type": "mp4",
					"url":  fmt.Sprintf("%s/%s/%s/g-720p.mp4", baseUrl, strconv.FormatUint(uint64(payload.ID), 10), info.VideoID),
				},
			},
		})
	}

	util.RenderPage(c, http.StatusOK, "UploadListPage", props)
}
