package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"testing"
)

var rawS3Event = []byte(`{
  "Records":[
    {
      "eventVersion":"2.0",
      "eventSource":"aws:s3",
      "awsRegion":"us-east-1",
      "eventTime":"2016-12-11T00:00:00.000Z",
      "eventName":"ObjectCreated:Put",
      "userIdentity":{
        "principalId":"A3MCB9FEJCFJSY"
      },
      "requestParameters":{
        "sourceIPAddress":"127.0.0.1"
      },
      "responseElements":{
        "x-amz-request-id":"3966C864F562A6A0",
        "x-amz-id-2":"2radsa8X4nKpba7KbgVurmc7rwe/SDoYLFid6MZKn18Nocpe3Ofwo5TJ+uJCnkf/"
      },
      "s3":{
        "s3SchemaVersion":"1.0",
        "configurationId":"Video Upload",
        "bucket":{
          "name":"serverless-video-upload",
          "ownerIdentity":{
            "principalId":"A3MCB9FEJCFJSY"
          },
          "arn":"arn:aws:s3:::serverless-video-upload"
        },
        "object":{
          "key":"my video.mp4",
          "size":2236480,
          "eTag":"ddb7a52094d2079a27ac44f83ca669e9",
          "sequencer": "005686091F4FFF1565"
        }
      }
    }
  ]
}`)

type TranscoderMock struct {
	createJobInvoked bool
}

func (t *TranscoderMock) CreateJob(input *elastictranscoder.CreateJobInput) (*elastictranscoder.CreateJobResponse, error) {
	t.createJobInvoked = true
	return &elastictranscoder.CreateJobResponse{}, nil
}

func TestInvocation(t *testing.T) {
	var s3Event events.S3Event
	err := json.Unmarshal(rawS3Event, &s3Event)
	if err != nil {
		t.Error(err)
	}

	var transcoderMock TranscoderMock
	err = createTranscoderJobOnNewS3Video(s3Event, &transcoderMock)
	if err != nil {
		t.Error(err)
	}
	if !transcoderMock.createJobInvoked {
		t.Error("CreateJob was never called")
	}
}
