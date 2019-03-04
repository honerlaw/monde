package controller

import (
	"github.com/gin-gonic/gin"
	"package/util"
	"net/http"
	"os"
	"package/middleware"
	"strconv"
	"github.com/satori/go.uuid"
	"time"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func UploadController(router *gin.Engine) {
	router.GET("/upload", renderUploadPage);
}

func renderUploadPage(c *gin.Context) {
	payload := middleware.GetAuthPayload(c);
	if payload == nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	id, _ := uuid.NewV4()
	userId := strconv.FormatInt(payload.ID, 10)
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

	policy := []byte(`{
		"expiration": "2020-12-01T12:00:00.000Z",
		"conditions": [
			{"acl": "`+acl+`"},
			{"bucket": "`+bucket+`"},
			{"x-amz-meta-user-id": "`+userId+`"},
			{"x-amz-algorithm": "`+algorithm+`"},
			{"x-amz-credential": "`+credential+`"},
			{"x-amz-date": "`+dateIso8601+`"},
			["starts-with", "$key", "process/`+userId+`/"]
		]
	}`)

	policyBase64 := base64.StdEncoding.EncodeToString(policy)

	// generate the signature
	dateKey := makeHmac([]byte("AWS4"+secretKey), []byte(shortTime))
	regionKey := makeHmac(dateKey, []byte(region))
	serviceKey := makeHmac(regionKey, []byte(service))
	credentialKey := makeHmac(serviceKey, []byte("aws4_request"))
	signatureHmac := makeHmac(credentialKey, []byte(policyBase64))
	signature := hex.EncodeToString(signatureHmac)

	props := gin.H{
		"authPayload": payload,
		"uploadBucketUrl": "http://" + os.Getenv("AWS_UPLOAD_BUCKET") + ".s3.amazonaws.com/",
		"uploadParams": gin.H{
			"acl":                acl,
			"key":                "process/" + userId + "/" + id.String(),
			"x-amz-meta-user-id": userId,
			"policy":             policyBase64,
			"x-amz-algorithm":    algorithm,
			"x-amz-credential":   credential,
			"x-amz-date":         dateIso8601,
			"x-amz-signature":    signature,
		},
	}

	util.RenderPage(c, http.StatusOK, "UploadPage", props)
}

func makeHmac(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}
