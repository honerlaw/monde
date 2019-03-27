package aws

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"time"
	"sync"
)

var s3Once sync.Once
var s3Instance *S3Service

type S3Service struct {
	client *s3.S3
}

func GetS3Service() (*S3Service) {
	s3Once.Do(func() {
		s3Instance = &S3Service{
			client: s3.New(Session),
		}
	})
	return s3Instance
}

func (service *S3Service) GetClient() (*s3.S3) {
	return service.client
}

func (service *S3Service) GetSignedUrl(bucket string, key string) (*string, error) {
	req, _ := service.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	url, presignErr := req.Presign(1 * time.Minute)

	if presignErr != nil {
		return nil, presignErr
	}

	return &url, nil;
}
