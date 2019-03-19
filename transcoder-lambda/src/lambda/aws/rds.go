package aws

import (
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"github.com/aws/aws-sdk-go/service/rds"
	"errors"
	aws2 "package/service/aws"
)

var _rdsClient *rds.RDS
var _rdsDataService *rdsdataservice.RDSDataService

func getRDSClient() (*rds.RDS) {
	if _rdsClient == nil {
		_rdsClient = rds.New(aws2.Session)
	}
	return _rdsClient
}

func getRDSDataService() (*rdsdataservice.RDSDataService) {
	if _rdsDataService == nil {
		_rdsDataService = rdsdataservice.New(aws2.Session)
	}
	return _rdsDataService
}

func getDbClusterArn() (*string, error) {
	client := getRDSClient()

	clusters, err := client.DescribeDBClusters(&rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(os.Getenv("CLUSTER_IDENTIFIER")),
	})

	if err != nil {
		return nil, err
	}

	if len(clusters.DBClusters) != 1 {
		return nil, errors.New("multiple db clusters found for given cluster identifier")
	}

	return clusters.DBClusters[0].DBClusterArn, nil
}

func RDSExecuteSQL(sql string) (*rdsdataservice.ExecuteSqlOutput, error) {
	secretArn, err := aws2.GetSMService().GetRDSSecretArn(os.Getenv("RDS_SECRET_NAME"))
	if err != nil {
		return nil, err
	}

	clusterArn, err := getDbClusterArn()
	if err != nil {
		return nil, err
	}

	output, err := getRDSDataService().ExecuteSql(&rdsdataservice.ExecuteSqlInput{
		AwsSecretStoreArn:      secretArn,
		DbClusterOrInstanceArn: clusterArn,
		Database:               aws.String(os.Getenv("DATABASE")),
		SqlStatements:          aws.String(sql),
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
