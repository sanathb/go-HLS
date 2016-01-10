package goAWS

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
)

const (

	//Preset ids
	//Create custom preset ids from the AWS console or use available ones for MP4 and webm format
	AWS_MP4_PRESET_ID  = ""
	AWS_WEBM_PRESET_ID = ""

	//These Preset ids used for HLS are already available in AWS
	AWS_HLS_AUDIO_PRESET_ID = "1351620000001-200071"
	AWS_HLS_0400K_PRESET_ID = "1351620000001-200050"
	AWS_HLS_0600K_PRESET_ID = "1351620000001-200040"
	AWS_HLS_1000K_PRESET_ID = "1351620000001-200030"
	AWS_HLS_1500K_PRESET_ID = "1351620000001-200020"
	AWS_HLS_2000K_PRESET_ID = "1351620000001-200010"

	//Sub folders for each HLS resolution
	AWS_HLS_AUDIO_FOLDER = "/hlsAudio/hlsAudio"
	AWS_HLS_0400k_FOLDER = "/hls0400k/hls0400k"
	AWS_HLS_0600K_FOLDER = "/hls0600k/hls0600k"
	AWS_HLS_1000K_FOLDER = "/hls1000k/hls1000k"
	AWS_HLS_1500K_FOLDER = "/hls1500k/hls1500k"
	AWS_HLS_2000K_FOLDER = "/hls2000k/hls2000k"

	AWS_HLS_FORMAT      = "HLSv3"
	AWS_SEGMET_DURATION = "10" //split into files having 10 sec duration
	AWS_REGION          = "us-east-1"

	AWS_VIDEO_INPUT           = "videos/"           //input video file location in S3 bucket
	AWS_TRANSCODE_PIPELINE_ID = ""                  //create a pipeline id from AWS console
	AWS_VIDEO_THUMBNAILS      = "video-thumbnails/" //folder in s3 bucket where the thumbnail has to be created
)

// To transcode video to HLS, mp4 and webm format
func TranscodeVideoHLS(filename string, extension string) {

	//Your panic handler
	//defer PanicHandler("TranscodeVideo")

	//export credentials at runtime
	//AWS_SECRET_ACCESS_KEY,AWS_ACCESS_KEY_ID

	creds := credentials.NewEnvCredentials()

	svc := elastictranscoder.New(session.New(), &aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: creds,
	})

	var (
		awsVideoOutputHLSAudio string = AWS_VIDEO_INPUT + filename + AWS_HLS_AUDIO_FOLDER
		awsVideoOutputHLS0400k string = AWS_VIDEO_INPUT + filename + AWS_HLS_0400k_FOLDER
		awsVideoOutputHLS0600k string = AWS_VIDEO_INPUT + filename + AWS_HLS_0600K_FOLDER
		awsVideoOutputHLS1000k string = AWS_VIDEO_INPUT + filename + AWS_HLS_1000K_FOLDER
		awsVideoOutputHLS1500k string = AWS_VIDEO_INPUT + filename + AWS_HLS_1500K_FOLDER
		awsVideoOutputHLS2000k string = AWS_VIDEO_INPUT + filename + AWS_HLS_2000K_FOLDER

		//Output folder for webm and mp4 format
		awsVideoOutputWebmMp4 string = AWS_VIDEO_INPUT + filename + "/" + filename
	)

	params := &elastictranscoder.CreateJobInput{
		Input: &elastictranscoder.JobInput{
			Key: aws.String(AWS_VIDEO_INPUT + filename + "." + extension),
		},
		PipelineId: aws.String(AWS_TRANSCODE_PIPELINE_ID),

		Outputs: []*elastictranscoder.CreateJobOutput{

			//mp4
			{
				//Thumbnail for the video is created here. Present this as a thumbnail for rest of the formats.
				Key:              aws.String(awsVideoOutputWebmMp4 + ".mp4"),
				PresetId:         aws.String(AWS_MP4_PRESET_ID),
				ThumbnailPattern: aws.String(AWS_VIDEO_THUMBNAILS + filename + "-{count}"),
			},
			//webm
			{
				Key:      aws.String(awsVideoOutputWebmMp4 + ".webm"),
				PresetId: aws.String(AWS_WEBM_PRESET_ID),
			},
			//HLS
			{
				Key:             aws.String(awsVideoOutputHLSAudio),
				PresetId:        aws.String(AWS_HLS_AUDIO_PRESET_ID),
				SegmentDuration: aws.String(AWS_SEGMET_DURATION),
			}, {
				Key:             aws.String(awsVideoOutputHLS0400k),
				PresetId:        aws.String(AWS_HLS_0400K_PRESET_ID),
				SegmentDuration: aws.String(AWS_SEGMET_DURATION),
			}, {
				Key:             aws.String(awsVideoOutputHLS0600k),
				PresetId:        aws.String(AWS_HLS_0600K_PRESET_ID),
				SegmentDuration: aws.String(AWS_SEGMET_DURATION),
			}, {
				Key:             aws.String(awsVideoOutputHLS1000k),
				PresetId:        aws.String(AWS_HLS_1000K_PRESET_ID),
				SegmentDuration: aws.String(AWS_SEGMET_DURATION),
			}, {
				Key:             aws.String(awsVideoOutputHLS1500k),
				PresetId:        aws.String(AWS_HLS_1500K_PRESET_ID),
				SegmentDuration: aws.String(AWS_SEGMET_DURATION),
			}, {
				Key:             aws.String(awsVideoOutputHLS2000k),
				PresetId:        aws.String(AWS_HLS_2000K_PRESET_ID),
				SegmentDuration: aws.String(AWS_SEGMET_DURATION),
			},
		},

		//Create master playlist
		Playlists: []*elastictranscoder.CreateJobPlaylist{
			{
				Format: aws.String(AWS_HLS_FORMAT),
				Name:   aws.String(AWS_VIDEO_INPUT + filename + "/" + filename),
				OutputKeys: []*string{
					aws.String(awsVideoOutputHLS2000k),
					aws.String(awsVideoOutputHLS1500k),
					aws.String(awsVideoOutputHLS1000k),
					aws.String(awsVideoOutputHLS0600k),
					aws.String(awsVideoOutputHLS0400k),
					aws.String(awsVideoOutputHLSAudio),
				},
			},
		},
	}

	svc.CreateJob(params)

}
