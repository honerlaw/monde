package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
	aws2 "lambda/aws"
	"lambda/util"
	"log"
)

func Handler(ctx context.Context, event events.S3Event) {
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

		/*_, err = aws2.ExecuteSQL("select * from monde.test")
		if err != nil {
			log.Print("Failed to execute sql", err)
			continue
		}*/
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	err = aws2.InitSession()
	if err != nil {
		panic(err)
	}

	lambda.Start(Handler)
}
