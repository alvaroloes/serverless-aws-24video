package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/url"
	"strings"
)

var awsSession = session.Must(session.NewSession())
var s3Manager = s3.New(awsSession, aws.NewConfig().WithRegion("us-east-1"))

func setTranscodedVideoPermissions(event events.SNSEvent) error {
	log.Printf("Received the following event:\n %+v", event)

	// Unmarshal the S3 event inside the SNS event
	var s3Event events.S3Event
	err := json.Unmarshal([]byte(event.Records[0].SNS.Message), &s3Event)
	if err != nil {
		return err
	}

	// Get the bucked and file name (replace spaces by '+' and unescape)
	bucketName := s3Event.Records[0].S3.Bucket.Name
	rawInputFile := s3Event.Records[0].S3.Object.Key
	inputFile, err := url.QueryUnescape(strings.Replace(rawInputFile, " ", "+", -1))
	if err != nil {
		return err
	}

	params := &s3.PutObjectAclInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(inputFile),
		ACL:    aws.String("public-read"),
	}

	res, err := s3Manager.PutObjectAcl(params)

	if err != nil {
		return err
	}

	log.Printf("Permissions changed successfully: %+v", res)

	return nil
}

func main() {
	lambda.Start(setTranscodedVideoPermissions)
}
