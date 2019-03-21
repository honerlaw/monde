package aws

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"sync"
)

var smOnce sync.Once
var smInstance *SMService

type SMService struct {
	client *secretsmanager.SecretsManager
}

type SMSecret struct {
	ARN   string
	Value string
}

func GetSMService() (*SMService) {
	smOnce.Do(func() {
		smInstance = &SMService{
			client: secretsmanager.New(Session),
		}
	})
	return smInstance
}

func (service *SMService) GetClient() (*secretsmanager.SecretsManager) {
	return service.client
}

func (service *SMService) GetSecret(secretName string) (*SMSecret, error) {
	output, err := service.client.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})

	if err != nil {
		return nil, err
	}

	return &SMSecret{
		ARN:   *output.ARN,
		Value: *output.SecretString,
	}, nil
}
