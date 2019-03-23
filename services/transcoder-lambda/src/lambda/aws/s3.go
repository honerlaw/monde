package aws

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"errors"
	"strconv"
	aws2 "server/service/aws"
)

type S3RecordMetadata struct {
	Bucket  string
	Key     string
	UserId  uint
	VideoId string
}

func GetS3RecordMetadata(bucket string, key string) (*S3RecordMetadata, error) {
	head, err := aws2.GetS3Service().GetClient().HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	});
	if err != nil {
		return nil, err
	}

	rawUserId, foundUserId := head.Metadata["User-Id"];
	rawVideoId, foundVideoId := head.Metadata["Video-Id"];
	if !foundUserId || !foundVideoId {
		return nil, errors.New("could not find required data in s3 user metadata")
	}

	userId, _ := strconv.ParseUint(*rawUserId, 10, 32)

	return &S3RecordMetadata{
		Bucket:  bucket,
		Key:     key,
		UserId:  uint(userId),
		VideoId: *rawVideoId,
	}, nil
}

