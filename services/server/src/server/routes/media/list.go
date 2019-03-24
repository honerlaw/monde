package media

import (
	"github.com/gin-gonic/gin"
	"server/util"
	"net/http"
	"server/middleware/auth"
	"server/service/media"
	"server/service/aws"
	"os"
	"fmt"
	"strconv"
)

func List(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY").(*auth.AuthPayload)
	props := gin.H{
		"authPayload": payload,
		"uploads": []gin.H{},
	}

	if payload != nil {

		infos, err := media.GetMediaInfoByUserId(payload.ID)
		if err != nil {
			// @todo redirect to 500 page or something
			return
		}

		for _, info := range *infos {
			status := aws.GetETService().GetJobStatus(info.JobID)

			baseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s", os.Getenv("AWS_PROCESSED_BUCKET"))
			// thumbBaseUrl := fmt.Sprintf("https://s3.amazonaws.com/%s/", os.Getenv("AWS_THUMBNAIL_BUCKET"))

			props["uploads"] = append(props["uploads"].([]gin.H), gin.H{
				"videoId": info.VideoID,
				"info": gin.H{
					"title": info.Title,
					"description": info.Description,
					"status": status,
				},
				"thumbs": []string{""}, // @todo figure out what this actually is
				"videos": []gin.H{
					{
						"type": "hls",
						"url": fmt.Sprintf("%s/%s/%s/playlist.m3u8", baseUrl, strconv.FormatUint(uint64(payload.ID), 10), info.VideoID),
					},
					{
						"type": "mp4",
						"url": fmt.Sprintf("%s/%s/%s/g-720p.mp4", baseUrl, strconv.FormatUint(uint64(payload.ID), 10), info.VideoID),
					},
				},
			})
		}
	}

	util.RenderPage(c, http.StatusOK, "UploadListPage", props)
}
