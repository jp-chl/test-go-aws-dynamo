package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jp-chl/test-go-aws-dynamo/src/db"

	"github.com/aws/aws-lambda-go/events"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

type DynamoDBService interface {
	GetItem(ctx context.Context, id string) (*db.Item, error)
}

func HandleRequest(ctx context.Context, request Request) (Response, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return handleError(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	dynamoService := db.NewDynamoDBService(client)
	return HandleRequestWithService(ctx, request, dynamoService)
}

func HandleRequestWithService(ctx context.Context, request Request, dynamoService DynamoDBService) (Response, error) {
	item, err := dynamoService.GetItem(ctx, request.PathParameters["id"])
	if err != nil {
		return handleError(err)
	}

	responseBody, err := json.Marshal(&item)
	if err != nil {
		return handleError(err)
	}

	response := Response{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		IsBase64Encoded: false,
	}

	return response, nil
}

func handleError(err error) (Response, error) {
	return Response{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, err
}
