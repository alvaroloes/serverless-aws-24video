#!/usr/bin/env bash

FUNCTION_NAME="set-permissions"
FUNCTION_ARN="arn:aws:lambda:us-east-1:785355572843:function:$FUNCTION_NAME"

echo "-> Starting deployment of '$FUNCTION_NAME'"

echo "-> Testing and vetting..."
go test -vet all

echo "-> Compiling..."
GOOS=linux GOARCH=amd64 go build -o ./deploy/${FUNCTION_NAME} main.go

echo "-> Zipping..."
zip -j ./deploy/${FUNCTION_NAME}.zip ./deploy/${FUNCTION_NAME}

echo "-> Deploying..."
aws lambda update-function-code \
  --function-name=${FUNCTION_ARN} \
  --zip-file=fileb://deploy/${FUNCTION_NAME}.zip

echo "-> Deployment finished."


