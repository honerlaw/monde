package service

import (
	"time"
	"strings"
	"regexp"
	"server/media/model"
	"errors"
	"log"
	"server/core/repository"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"strconv"
	"os"
	"encoding/base64"
	"encoding/hex"
	"server/core/util"
	"server/user/middleware"
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

func Update(req UpdateRequest) (error) {
	tx := repository.DB.Begin()

	info, err := GetMediaInfoByVideoID(req.VideoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	info.Title = req.Title
	info.Description = req.Description

	err = Save(info)
	if err != nil {
		tx.Rollback()
		return err
	}

	var hashtags []model.Hashtag
	regex, _ := regexp.Compile("^#\\w+")
	tags := strings.Split(req.Hashtags, " ")
	for _, tag := range tags {
		trimmedTag := strings.TrimSpace(tag)
		if regex.MatchString(trimmedTag) && len(trimmedTag) < 75 && len(trimmedTag) > 1 {
			// fetch the existing hashtag if it exists
			hashtag, err := GetHashTag(trimmedTag)
			if hashtag == nil || err != nil {
				hashtag = &model.Hashtag{
					Tag: trimmedTag,
				}
			}

			hashtags = append(hashtags, *hashtag)
		}
	}

	assoc := tx.Model(info).Association("Hashtags").Replace(hashtags)
	if assoc.Error != nil {
		tx.Rollback()
		log.Print(assoc.Error)
		return errors.New("failed to update media")
	}

	tx.Commit()

	return nil
}

func TogglePublish(req PublishRequest) (error) {
	info, err := GetMediaInfoByVideoID(req.VideoID)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(info.Description)) == 0 {
		return errors.New("a description is required to publish videos")
	}

	info.Published = !info.Published
	if info.Published {
		info.PublishedDate = time.Now()
	}

	err = Save(info)
	if err != nil {
		return err
	}

	return nil
}

func GetUploadFormProps(payload *middleware.AuthPayload) (*gin.H) {
	id, _ := uuid.NewV4()
	userId := strconv.FormatUint(uint64(payload.ID), 10)
	bucket := os.Getenv("AWS_UPLOAD_BUCKET")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	algorithm := "AWS4-HMAC-SHA256"
	acl := "public-read-write"
	currentTime := time.Now()
	dateIso8601 := currentTime.Format("20060102T150405Z")
	shortTime := currentTime.Format("20060102")
	service := "s3"
	credential := accessKey + "/" + shortTime + "/" + region + "/" + service + "/aws4_request"
	redirect := os.Getenv("UPLOAD_REDIRECT_DOMAIN") + "/media/list"

	policy := []byte(`{
		"expiration": "2020-12-01T12:00:00.000Z",
		"conditions": [
			{"acl": "` + acl + `"},
			{"bucket": "` + bucket + `"},
			{"success_action_redirect": "` + redirect + `"},
			{"x-amz-meta-user-id": "` + userId + `"},
			{"x-amz-meta-video-id": "` + id.String() + `"},
			{"x-amz-algorithm": "` + algorithm + `"},
			{"x-amz-credential": "` + credential + `"},
			{"x-amz-date": "` + dateIso8601 + `"},
			["starts-with", "$key", "` + userId + `/"]
		]
	}`)

	policyBase64 := base64.StdEncoding.EncodeToString(policy)

	// generate the signature
	dateKey := util.MakeHmac([]byte("AWS4"+secretKey), []byte(shortTime))
	regionKey := util.MakeHmac(dateKey, []byte(region))
	serviceKey := util.MakeHmac(regionKey, []byte(service))
	credentialKey := util.MakeHmac(serviceKey, []byte("aws4_request"))
	signatureHmac := util.MakeHmac(credentialKey, []byte(policyBase64))
	signature := hex.EncodeToString(signatureHmac)

	return &gin.H{
		"uploadBucketUrl": "http://" + os.Getenv("AWS_UPLOAD_BUCKET") + ".s3.amazonaws.com/",
		"uploadParams": gin.H{
			"acl":                     acl,
			"key":                     userId + "/" + id.String(),
			"success_action_redirect": redirect,
			"x-amz-meta-user-id":      userId,
			"x-amz-meta-video-id":     id.String(),
			"policy":                  policyBase64,
			"x-amz-algorithm":         algorithm,
			"x-amz-credential":        credential,
			"x-amz-date":              dateIso8601,
			"x-amz-signature":         signature,
		},
	}
}
