package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jp-chl/test-go-aws-dynamo/src/db"

	"github.com/aws/aws-lambda-go/events"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

func HandleRequest(ctx context.Context, request Request) (Response, error) {
	item, err := db.GetItem(ctx, request.PathParameters["id"])
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
