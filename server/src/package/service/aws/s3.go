package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

var s3Client *s3.S3

func GetS3Client() (*s3.S3) {
	if s3Client != nil {
		return s3Client
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewEnvCredentials(),
	})

	if err != nil {
		panic(err)
	}

	s3Client = s3.New(sess)

	return s3Client
}
