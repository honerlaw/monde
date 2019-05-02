package aws

import (
	"sync"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"github.com/labstack/gommon/log"
)

var sesOnce sync.Once
var sesInstance *SESService

type SESService struct {
	client *ses.SES
}

func GetSESService() (*SESService) {
	sesOnce.Do(func() {
		sesInstance = &SESService{
			client: ses.New(Session),
		}
	})
	return sesInstance
}

func (service *SESService) GetClient() (*ses.SES) {
	return service.client
}

func (service *SESService) SendEmail(to string, subject string, message string) (error) {
	_, err := service.client.SendEmail(&ses.SendEmailInput{
		Source: aws.String(os.Getenv("EMAIL")),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data: aws.String(subject),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data: aws.String(message),
				},
			},
		},
	})

	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
