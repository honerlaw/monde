package aws

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"strconv"
)

type S3RecordMetadata struct {
	UserId int64
}

var _s3Client *s3.S3

func getS3Client() (*s3.S3) {
	if _s3Client == nil {
		_s3Client = s3.New(Session)
	}
	return _s3Client
}

func GetS3RecordMetadata(bucket string, key string) (*S3RecordMetadata, error) {
	head, err := getS3Client().HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	});

	if err != nil {
		return nil, err
	}

	userId, _ := strconv.ParseInt(*head.Metadata["User-Id"], 10, 64)

	return &S3RecordMetadata{
		UserId: userId,
	}, nil
}
