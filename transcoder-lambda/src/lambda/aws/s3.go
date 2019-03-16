package aws

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"strconv"
	"errors"
	"time"
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

func GetSignedS3Url(bucket string, key string) (*string, error) {
	req, _ := getS3Client().GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	url, presignErr := req.Presign(1 * time.Minute)

	if presignErr != nil {
		return nil, presignErr
	}

	return &url, nil;
}

func GetS3RecordMetadata(bucket string, key string) (*S3RecordMetadata, error) {
	head, err := getS3Client().HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	});

	if err != nil {
		return nil, err
	}

	rawUserId, foundUserId := head.Metadata["User-Id"];

	if !foundUserId {
		return nil, errors.New("could not find user id in user metadata")
	}

	userId, _ := strconv.ParseInt(*rawUserId, 10, 64)

	return &S3RecordMetadata{
		UserId: userId,
	}, nil
}
