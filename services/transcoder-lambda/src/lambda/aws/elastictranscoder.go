package aws

import (
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/aws/aws-sdk-go/aws"
	"strconv"
	"strings"
	aws2 "package/service/aws"
	"os"
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
	{
		Id: "1351620000001-000010",
		Name: "g-720p.mp4",
		FullName: "Generic 720p",
	},
}

func CreateElasticTranscoderJob(metadata *S3RecordMetadata) (*elastictranscoder.Job, error) {
	pipelineId, err := aws2.GetETService().GetPipelineID(os.Getenv("TRANSCODER_PIPELINE_NAME"))
	if err != nil {
		return nil, err
	}

	var playlistOutputKeys []*string
	var outputs []*elastictranscoder.CreateJobOutput

	// build the outputs from the presets defined above, ideally we would fetch this info from the database or
	// somewhere else instead of have them hardcoded...
	for _, preset := range presets {
		output := &elastictranscoder.CreateJobOutput{
			Key:              aws.String(preset.Name),
			PresetId:         aws.String(preset.Id),
			ThumbnailPattern: aws.String(preset.Name + "-{count}"),
			Rotate:           aws.String("0"),
		}

		// we only want to add the hls files to be added to the playlist
		if strings.HasPrefix(preset.Name, "hls-") {
			output.Key = aws.String("video/" + preset.Name)
			playlistOutputKeys = append(playlistOutputKeys, output.Key)
			output.SegmentDuration = aws.String("5")
		}

		outputs = append(outputs, output)
	}

	resp, err := aws2.GetETService().GetClient().CreateJob(&elastictranscoder.CreateJobInput{
		PipelineId: pipelineId,
		Input: &elastictranscoder.JobInput{
			Key:         aws.String(metadata.Key),
			AspectRatio: aws.String("auto"),
			Container:   aws.String("auto"),
			FrameRate:   aws.String("auto"),
			Interlaced:  aws.String("auto"),
			Resolution:  aws.String("auto"),
		},
		OutputKeyPrefix: aws.String(strconv.FormatUint(uint64(metadata.UserId), 10) + "/" + metadata.VideoId + "/"),
		Outputs:         outputs,
		Playlists: []*elastictranscoder.CreateJobPlaylist{
			{
				Format:     aws.String("HLSv3"),
				Name:       aws.String("playlist"),
				OutputKeys: playlistOutputKeys,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return resp.Job, nil
}
