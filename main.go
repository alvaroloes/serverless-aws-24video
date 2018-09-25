package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"log"
	"net/url"
	"path"
	"strings"
)

func createTranscoderJobOnNewS3Video(ctx context.Context, event events.S3Event) error {
	log.Printf("Received the following event:\n %+v", event)

	// Get and unescape the input file (replace spaces by '+' prior unescaping)
	rawInputFile := event.Records[0].S3.Object.Key
	inputFile, err := url.QueryUnescape(strings.Replace(rawInputFile, " ", "+", -1))
	if err != nil {
		return err
	}

	// Prepare the output file name
	outputFile := strings.TrimSuffix(inputFile, path.Ext(inputFile))

	// Prepare the parameters for the job
	params := &elastictranscoder.CreateJobInput{
		PipelineId: aws.String("1537526232182-ip27yr"),
		Input: &elastictranscoder.JobInput{
			Key: aws.String(inputFile),
		},
		Outputs: []*elastictranscoder.CreateJobOutput{
			{
				Key:      aws.String(outputFile + "-1080p" + ".mp4"),
				PresetId: aws.String("1351620000001-000001"), //Generic 1080p
			},
			{
				Key:      aws.String(outputFile + "-720p" + ".mp4"),
				PresetId: aws.String("1351620000001-000010"), //Generic 720p
			},
			{
				Key:      aws.String(outputFile + "-web-720p" + ".mp4"),
				PresetId: aws.String("1351620000001-100070"), //Web Friendly 720p
			},
		},
	}

	// Create the transcoder
	awsSession, err := session.NewSession()
	if err != nil {
		return err
	}
	var transcoder = elastictranscoder.New(awsSession, aws.NewConfig().WithRegion("us-east-1"))

	job, err := transcoder.CreateJob(params)
	if err != nil {
		return err
	}

	log.Printf("Job created successfully: %+v", job)

	return nil
}

func main() {
	lambda.Start(createTranscoderJobOnNewS3Video)
}
