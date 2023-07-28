package main

import (
	"github.com/jp-chl/test-go-aws-dynamo/src/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.HandleRequest)
}
