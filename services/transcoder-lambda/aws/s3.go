package aws

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"errors"
	aws2 "services/server/core/service/aws"
)

type S3RecordMetadata struct {
	Bucket string
	Key    string
	UserID string
	ID     string
}

func GetS3RecordMetadata(bucket string, key string) (*S3RecordMetadata, error) {
	head, err := aws2.GetS3Service().GetClient().HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	});
	if err != nil {
		return nil, err
	}

	rawUserID, foundUserId := head.Metadata["User-Id"];
	rawID, foundVideoId := head.Metadata["Id"];
	if !foundUserId || !foundVideoId {
		return nil, errors.New("could not find required data in s3 user metadata")
	}

	return &S3RecordMetadata{
		Bucket: bucket,
		Key:    key,
		UserID: *rawUserID,
		ID:     *rawID,
	}, nil
}
