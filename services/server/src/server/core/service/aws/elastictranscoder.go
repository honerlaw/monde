package aws

import (
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"errors"
	"sync"
	"github.com/aws/aws-sdk-go/aws"
	"log"
)

const (
	ETJobStatusCanceled    = "Canceled"
	ETJobStatusProgressing = "Progressing"
	ETJobStatusComplete    = "Complete"
	ETJobStatusSubmitted   = "Submitted"
	ETJobStatusError       = "Error"
)

var etOnce sync.Once
var etInstance *ETService

type ETService struct {
	client *elastictranscoder.ElasticTranscoder
}

func GetETService() (*ETService) {
	etOnce.Do(func() {
		etInstance = &ETService{
			client: elastictranscoder.New(Session),
		}
	})
	return etInstance
}

func (service *ETService) GetClient() (*elastictranscoder.ElasticTranscoder) {
	return service.client
}

func (service *ETService) GetPipelineID(name string) (*string, error) {
	var pageToken *string

	// continue searching until we either error out, or run out of pages
	for true {

		// fetch a list of pipelines
		output, err := service.client.ListPipelines(&elastictranscoder.ListPipelinesInput{
			PageToken: pageToken,
		})

		// if the request errored, just short circuit with the error
		if err != nil {
			return nil, err
		}

		// otherwise see if the pipeline we are looking for exists
		for _, pipeline := range output.Pipelines {

			// if it does, short circuit and return the id
			if *pipeline.Name == name {
				return pipeline.Id, nil
			}
		}

		// if there is a page token, set it for the next request
		pageToken = output.NextPageToken
		if pageToken == nil {
			break
		}
	}

	return nil, errors.New("could not find pipeline with name: " + name)
}

func (service *ETService) GetJobStatus(jobId string) (string) {
	if len(jobId) == 0 {
		return ETJobStatusError
	}

	output, err := service.client.ReadJob(&elastictranscoder.ReadJobInput{
		Id: aws.String(jobId),
	})

	if err != nil {
		log.Print("failed to fetch job status", err)
		return ETJobStatusError
	}

	// we return a enum instead of a string, so we can more easily just compare shit
	return *output.Job.Status;
}
