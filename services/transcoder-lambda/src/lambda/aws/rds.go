package aws

import (
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	aws2 "package/service/aws"
	"sync"
)

var rdsdsSync sync.Once
var rdsdsInstance *RDSDService

type RDSDService struct {
	client *rdsdataservice.RDSDataService
}

func GetRDSDService() (*RDSDService) {
	rdsdsSync.Do(func() {
		rdsdsInstance = &RDSDService{
			client: rdsdataservice.New(aws2.Session),
		}
	})
	return rdsdsInstance
}

func (service *RDSDService) ExecuteSQL(sql string) (*rdsdataservice.ExecuteSqlOutput, error) {
	secret, err := aws2.GetSMService().GetSecret(os.Getenv("DB_SECRET_NAME"))
	if err != nil {
		return nil, err
	}

	cluster, err := aws2.GetRDSService().GetCluster()
	if err != nil {
		return nil, err
	}

	output, err := service.client.ExecuteSql(&rdsdataservice.ExecuteSqlInput{
		AwsSecretStoreArn:      aws.String(secret.ARN),
		DbClusterOrInstanceArn: aws.String(cluster.ARN),
		Database:               aws.String(os.Getenv("DB_NAME")),
		SqlStatements:          aws.String(sql),
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
