package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jp-chl/test-go-aws-dynamo/src/db"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

type DynamoDBService interface {
	GetItem(ctx context.Context, id string) (*db.Item, error)
}

// MarshalWrapper wraps attributevalue.Marshal to match the AttributeMarshaler interface
type MarshalWrapper struct{}

func (mw *MarshalWrapper) Marshal(in interface{}) (types.AttributeValue, error) {
	return attributevalue.Marshal(in)
}

// UnmarshalMapWrapper wraps attributevalue.UnmarshalMap to match the AttributeUnmarshaler interface
type UnmarshalMapWrapper struct{}

func (umw *UnmarshalMapWrapper) UnmarshalMap(m map[string]types.AttributeValue, out interface{}) error {
	return attributevalue.UnmarshalMap(m, out)
}

func HandleRequest(ctx context.Context, request Request) (Response, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return handleError(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	dynamoService := db.NewDynamoDBService(client, &MarshalWrapper{}, &UnmarshalMapWrapper{})
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
