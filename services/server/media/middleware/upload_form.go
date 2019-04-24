package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
	"time"
	"encoding/base64"
	"services/server/core/util"
	"encoding/hex"
	"github.com/satori/go.uuid"
	"services/server/user/middleware"
	"services/server/core/render"
	"services/server/user/service"
)

func UploadFormMiddleware(channelService *service.ChannelService) (gin.HandlerFunc) {
	return func(c *gin.Context) {
		c.Next()

		payload := c.MustGet("JWT_AUTH_PAYLOAD")
		if payload != nil {

			meta, metaExists := c.Get("react-meta")
			if meta != nil && metaExists {
				meta.(render.ReactMeta).Props["uploadForm"] = getUploadFormProps(payload.(*middleware.AuthPayload), channelService)
			}
		}
	}
}

func getUploadFormProps(payload *middleware.AuthPayload, channelService *service.ChannelService) (*gin.H) {
	channel, _ := channelService.GetByUserID(payload.ID, nil)

	id := uuid.NewV4()
	channelId := channel.ID
	bucket := os.Getenv("AWS_UPLOAD_BUCKET")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	algorithm := "AWS4-HMAC-SHA256"
	acl := "public-read-write"
	currentTime := time.Now()
	dateIso8601 := currentTime.Format("20060102T150405Z")
	shortTime := currentTime.Format("20060102")
	serviceName := "s3"
	credential := accessKey + "/" + shortTime + "/" + region + "/" + serviceName + "/aws4_request"
	redirect := os.Getenv("UPLOAD_REDIRECT_DOMAIN") + "/media/list"

	policy := []byte(`{
		"expiration": "2020-12-01T12:00:00.000Z",
		"conditions": [
			{"acl": "` + acl + `"},
			{"bucket": "` + bucket + `"},
			{"success_action_redirect": "` + redirect + `"},
			{"x-amz-meta-channel-id": "` + channelId + `"},
			{"x-amz-meta-video-id": "` + id.String() + `"},
			{"x-amz-algorithm": "` + algorithm + `"},
			{"x-amz-credential": "` + credential + `"},
			{"x-amz-date": "` + dateIso8601 + `"},
			["starts-with", "$key", "` + channelId + `/"]
		]
	}`)

	policyBase64 := base64.StdEncoding.EncodeToString(policy)

	// generate the signature
	dateKey := util.MakeHmac([]byte("AWS4"+secretKey), []byte(shortTime))
	regionKey := util.MakeHmac(dateKey, []byte(region))
	serviceKey := util.MakeHmac(regionKey, []byte(serviceName))
	credentialKey := util.MakeHmac(serviceKey, []byte("aws4_request"))
	signatureHmac := util.MakeHmac(credentialKey, []byte(policyBase64))
	signature := hex.EncodeToString(signatureHmac)

	return &gin.H{
		"uploadBucketUrl": "http://" + os.Getenv("AWS_UPLOAD_BUCKET") + ".s3.amazonaws.com/",
		"uploadParams": gin.H{
			"acl":                     acl,
			"key":                     channelId + "/" + id.String(),
			"success_action_redirect": redirect,
			"x-amz-meta-channel-id":   channelId,
			"x-amz-meta-video-id":     id.String(),
			"policy":                  policyBase64,
			"x-amz-algorithm":         algorithm,
			"x-amz-credential":        credential,
			"x-amz-date":              dateIso8601,
			"x-amz-signature":         signature,
		},
	}
}
