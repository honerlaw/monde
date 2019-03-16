package aws

import (
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"errors"
	"strconv"
)

type VideoPreset struct {
	Id       string
	Name     string
	FullName string
}

var presets = []*VideoPreset{
	{
		Id:       "1351620000001-200015",
		Name:     "hls-v-2m",
		FullName: "HLS Video - 2M",
	},
	{
		Id:       "1351620000001-200025",
		Name:     "hls-v-1-5m",
		FullName: "HLS Video - 1.5M",
	},
	{
		Id:       "1351620000001-200035",
		Name:     "hls-v-1m",
		FullName: "HLS Video - 1M",
	},
	{
		Id:       "1351620000001-200045",
		Name:     "hls-v-600k",
		FullName: "HLS Video - 600k",
	},
	{
		Id:       "1351620000001-200055",
		Name:     "hls-v-400k",
		FullName: "HLS Video - 400k",
	},
}

var _etClient *elastictranscoder.ElasticTranscoder

func getETClient() (*elastictranscoder.ElasticTranscoder) {
	if _etClient == nil {
		_etClient = elastictranscoder.New(Session)
	}
	return _etClient
}

func getPipelineID() (*string, error) {
	name := os.Getenv("TRANSCODER_PIPELINE_NAME")

	var pageToken *string

	// continue searching until we either error out, or run out of pages
	for true {

		// fetch a list of pipelines
		output, err := getETClient().ListPipelines(&elastictranscoder.ListPipelinesInput{
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

func CreateElasticTranscoderJob(metadata *S3RecordMetadata) (*elastictranscoder.Job, error) {
	pipelineId, err := getPipelineID()
	if err != nil {
		return nil, err
	}

	var outputKeys []*string
	var outputs []*elastictranscoder.CreateJobOutput

	// build the outputs from the presets defined above, ideally we would fetch this info from the database or
	// somewhere else
	for _, preset := range presets {
		key := aws.String("video/" + preset.Name)

		outputKeys = append(outputKeys, key)

		outputs = append(outputs, &elastictranscoder.CreateJobOutput{
			Key:              key,
			PresetId:         aws.String(preset.Id),
			ThumbnailPattern: aws.String(preset.Name + "-{count}"),
			Rotate:           aws.String("0"),
			SegmentDuration:  aws.String("5"),
		})
	}

	resp, err := getETClient().CreateJob(&elastictranscoder.CreateJobInput{
		PipelineId: pipelineId,
		Input: &elastictranscoder.JobInput{
			Key:         aws.String(metadata.Key),
			AspectRatio: aws.String("auto"),
			Container:   aws.String("auto"),
			FrameRate:   aws.String("auto"),
			Interlaced:  aws.String("auto"),
			Resolution:  aws.String("auto"),
		},
		OutputKeyPrefix: aws.String(strconv.FormatInt(metadata.UserId, 10) + "/" + metadata.VideoId + "/"),
		Outputs:         outputs,
		Playlists: []*elastictranscoder.CreateJobPlaylist{
			{
				Format:     aws.String("HLSv3"),
				Name:       aws.String("playlist.m3u8"),
				OutputKeys: outputKeys,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return resp.Job, nil
}
