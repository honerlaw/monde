package aws

import (
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"github.com/aws/aws-sdk-go/service/rds"
	"errors"
	"package/model/media"
	"fmt"
)

var _rdsClient *rds.RDS
var _rdsDataService *rdsdataservice.RDSDataService

func getRDSClient() (*rds.RDS) {
	if _rdsClient == nil {
		_rdsClient = rds.New(Session)
	}
	return _rdsClient
}

func getRDSDataService() (*rdsdataservice.RDSDataService) {
	if _rdsDataService == nil {
		_rdsDataService = rdsdataservice.New(Session)
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

func execute(sql string) (*rdsdataservice.ExecuteSqlOutput, error) {
	secretArn, err := GetRDSSecretArn()
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

// this monstrosity should insert all of the media info properly
// We can't really use transactions because of the wonderful aurora data api limitations...
// We shouldn't need to sworry about escaping things because we generated all the data...
// HOWEVER, we should probably still try and escape crap just in case mediainfo gets some weird results back
// @todo escape things to prevent sql injections / unintended side effects
func Insert(mediainfo *media.MediaInfo) (error) {
	_, err := execute(fmt.Sprintf(
		`INSERT INTO media_info (created_at, user_id, job_id, video_id) VALUES (NOW(), %d, '%s', '%s')`,
		mediainfo.UserID, mediainfo.JobID, mediainfo.VideoID,
	))
	if err != nil {
		return err
	}

	resp, err := execute(fmt.Sprintf(
		"SELECT id FROM media_info WHERE user_id = %d AND job_id = '%s' and video_id = '%s'",
		mediainfo.UserID, mediainfo.JobID, mediainfo.VideoID,
	))
	if err != nil {
		return err
	}

	mediaInfoId := *resp.SqlStatementResults[0].ResultFrame.Records[0].Values[0].IntValue;
	for _, media := range mediainfo.Medias {
		_, err = execute(fmt.Sprintf(
			`INSERT INTO media (created_at, media_info_id) VALUES (NOW(), %d)`,
			mediaInfoId,
		))
		if err != nil {
			return err
		}

		resp, err := execute(fmt.Sprintf(
			"SELECT id FROM media WHERE media_info_id = %d ORDER BY created_at DESC LIMIT 1",
			mediaInfoId,
		))
		if err != nil {
			return err
		}
		mediaId := *resp.SqlStatementResults[0].ResultFrame.Records[0].Values[0].IntValue

		for _, track := range media.Tracks {
			_, err = execute(fmt.Sprintf(
				`INSERT INTO track (created_at, media_id, type, duration, width, height, format, encoded_date, video_count, data_size, file_size) VALUES (NOW(), %d, '%s', %f, %d, %d, '%s', '%s', '%s', %d, %d)`,
				mediaId, track.Type, track.Duration, track.Width, track.Height, track.Format, track.Encoded_Date, track.VideoCount, track.DataSize, track.FileSize,
			))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
