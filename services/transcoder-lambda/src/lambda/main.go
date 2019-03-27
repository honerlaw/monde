package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
	"lambda/util"
	"log"
	aws2 "lambda/aws"
	"server/core/service/aws"
	"server/core/repository"
)

func Handler(ctx context.Context, event events.S3Event) {
	// @todo we shouldn't need to close everytime, but for now we will
	defer repository.Connect().Close()

	for _, record := range event.Records {
		metadata, err := aws2.GetS3RecordMetadata(record.S3.Bucket.Name, record.S3.Object.Key)
		if err != nil {
			log.Print("Failed to get s3 metadata", err)
			continue
		}

		mediainfo, err := util.GetMediaInfo(metadata)
		if err != nil {
			log.Print("Failed to get media info from file", err)
			continue;
		}

		err = util.ValidateMediaInfo(mediainfo)
		if err != nil {
			log.Print("Invalid media file", err)
			continue;
		}

		job, err := aws2.CreateElasticTranscoderJob(metadata)
		if err != nil {
			log.Print("Failed to trigger elastic transcoder job", err)
			continue;
		}

		// set the job id so we can look it up later if we need to
		mediainfo.JobID = *job.Id

		err = util.Insert(mediainfo)
		if err != nil {
			log.Print("Failed to insert media and job informations", err)
			continue
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	err = aws.InitSession()
	if err != nil {
		panic(err)
	}

	lambda.Start(Handler)
}
