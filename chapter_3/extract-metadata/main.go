package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var awsSession = session.Must(session.NewSession())
var s3Manager = s3.New(awsSession, aws.NewConfig().WithRegion("us-east-1"))

func main() {
	lambda.Start(mainHandler)
}

func mainHandler(event events.SNSEvent) error {
	log.Printf("Received the following event:\n %+v", event)

	// Unmarshal the S3 event inside the SNS event
	var s3Event events.S3Event
	err := json.Unmarshal([]byte(event.Records[0].SNS.Message), &s3Event)
	if err != nil {
		return err
	}

	// Get the bucked and file name (replace spaces by '+' and unescape)
	bucketName := s3Event.Records[0].S3.Bucket.Name
	rawInputFileName := s3Event.Records[0].S3.Object.Key
	inputFileName, err := url.QueryUnescape(strings.Replace(rawInputFileName, " ", "+", -1))
	if err != nil {
		return err
	}

	// Save the video to the filesystem
	videoFilename := "/tmp/" + filepath.Base(inputFileName)
	err = saveS3FileToFilesystem(bucketName, inputFileName, videoFilename)
	if err != nil {
		return err
	}

	err = extractMetadata(videoFilename)

	return nil
}

func saveS3FileToFilesystem(bucketName, inputFileName, outputFilename string) error {
	log.Println("Saving S3 file to filesystem")

	// Get the S3 file input stream
	res, err := s3Manager.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(inputFileName),
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Create the file
	localFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// Copy from S3 to our file
	_, err = io.Copy(localFile, res.Body)
	return err
}

func extractMetadata(videFilename string) error {
	return nil
}
