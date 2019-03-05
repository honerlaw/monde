package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"fmt"
	"github.com/joho/godotenv"
	aws2 "lambda/aws"
)

// @todo create a transcoder job and ship it off
// @todo insert tracking metadata to know the status of the transcoder job
func Handler(ctx context.Context, event events.S3Event) {
	for _, record := range event.Records {
		_, err := aws2.GetS3RecordMetadata(record.S3.Bucket.Name, record.S3.Object.Key)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}
		_, err = aws2.ExecuteSQL("select * from table")
		if err != nil {
			fmt.Print(err.Error())
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
