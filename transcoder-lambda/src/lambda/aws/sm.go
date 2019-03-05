package aws

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"os"
)

var _smClient *secretsmanager.SecretsManager

func getSMClient() (*secretsmanager.SecretsManager) {
	if  _smClient == nil {
		_smClient = secretsmanager.New(Session)
	}
	return _smClient
}

func GetRDSSecretArn() (*string, error) {
	client := getSMClient()

	output, err := client.ListSecretVersionIds(&secretsmanager.ListSecretVersionIdsInput{
		SecretId: aws.String(os.Getenv("RDS_SECRET_NAME")),
	})

	if err != nil {
		return nil, err
	}

	return output.ARN, nil
}
