package aws

import (
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"strings"
)

var s3Client *s3.S3

func getS3Client() (*s3.S3) {
	if s3Client != nil {
		return s3Client
	}

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewEnvCredentials(),
	})

	if err != nil {
		panic(err)
	}

	// Create S3 service client
	s3Client = s3.New(sess)

	return s3Client
}

func S3GetPresignedUrl(bucket string, key string, metadata map[string]string, minutes time.Duration) (string, error) {
	s3Client := getS3Client()

	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		ACL:        aws.String("public-read-write"),
		Bucket:     aws.String(bucket),
		Key:        aws.String(key),
		Metadata:   aws.StringMap(metadata),
		Body:       nil,
		ContentMD5: nil,
	})

	u, err := req.Presign(minutes * time.Minute)

	u = strings.Replace(u, "%", "%%", -1)

	return u, err
}
