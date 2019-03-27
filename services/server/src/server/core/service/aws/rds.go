package aws

import (
	"sync"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"errors"
)

var rdsOnce sync.Once
var rdsInstance *RDSService

type RDSService struct {
	client *rds.RDS
}

type RDSCluster struct {
	ARN      string
	Endpoint string
}

func GetRDSService() (*RDSService) {
	rdsOnce.Do(func() {
		rdsInstance = &RDSService{
			client: rds.New(Session),
		}
	})
	return rdsInstance
}

func (service *RDSService) GetClient() (*rds.RDS) {
	return service.client
}

func (service *RDSService) GetCluster() (*RDSCluster, error) {
	clusters, err := service.client.DescribeDBClusters(&rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(os.Getenv("DB_CLUSTER_IDENTIFIER")),
	})

	if err != nil {
		return nil, err
	}

	if len(clusters.DBClusters) != 1 {
		return nil, errors.New("multiple db clusters found for given cluster identifier")
	}

	cluster := clusters.DBClusters[0]

	return &RDSCluster{
		ARN:      *cluster.DBClusterArn,
		Endpoint: *cluster.Endpoint,
	}, nil
}
