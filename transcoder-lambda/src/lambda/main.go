package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"fmt"
	"github.com/joho/godotenv"
	aws2 "lambda/aws"
	"lambda/util"
)

func Handler(ctx context.Context, event events.S3Event) {
	for _, record := range event.Records {
		_, err := aws2.GetS3RecordMetadata(record.S3.Bucket.Name, record.S3.Object.Key)
		if err != nil {
			fmt.Print("Failed to get s3 metadata", err)
			continue
		}

		_, err = util.GetMediaInfo(record.S3.Bucket.Name, record.S3.Object.Key)
		if err != nil {
			fmt.Print("Failed to get media info from file", err)
			continue;
		}

		_, err = aws2.ExecuteSQL("select * from monde.test")
		if err != nil {
			fmt.Print("Failed to execute sql", err)
			continue
		}
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
