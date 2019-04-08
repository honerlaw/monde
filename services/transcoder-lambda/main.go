package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
	"log"
	aws2 "services/transcoder-lambda/aws"
	"services/server/core/service/aws"
	"services/server/core/repository"
	"services/transcoder-lambda/service"
)

func Handler(ctx context.Context, event events.S3Event) {
	repository.GetRepository() // basically initializes the db connnection

	for _, record := range event.Records {
		metadata, err := aws2.GetS3RecordMetadata(record.S3.Bucket.Name, record.S3.Object.Key)
		if err != nil {
			log.Print("Failed to get s3 metadata", err)
			continue
		}

		mediainfo, err := service.GetMediaInfo(metadata)
		if err != nil {
			log.Print("Failed to get media info from file", err)
			service.RetryInsert(metadata, mediainfo, nil, 5)
			continue;
		}

		err = service.ValidateMediaInfo(mediainfo)
		if err != nil {
			log.Print("Invalid media file", err)
			service.RetryInsert(metadata, mediainfo, nil, 5)
			continue;
		}

		job, err := aws2.CreateElasticTranscoderJob(metadata)
		if err != nil {
			log.Print("Failed to trigger elastic transcoder job", err)
			service.RetryInsert(metadata, mediainfo, nil, 5)
			continue;
		}

		service.RetryInsert(metadata, mediainfo, job, 5)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = aws.InitSession()
	if err != nil {
		log.Fatal(err)
	}

	lambda.Start(Handler)
}
