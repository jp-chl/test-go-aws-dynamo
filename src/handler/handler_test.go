package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/jp-chl/test-go-aws-dynamo/src/db"
	"github.com/stretchr/testify/assert"
)

type mockDynamoDBService struct {
	GetItemFunc func(ctx context.Context, id string) (*db.Item, error)
}

func (m *mockDynamoDBService) GetItem(ctx context.Context, id string) (*db.Item, error) {
	return m.GetItemFunc(ctx, id)
}

func TestHandleRequest(t *testing.T) {
	ctx := context.TODO()
	request := Request{PathParameters: map[string]string{"id": "testID"}}

	// Success case
	mockService := &mockDynamoDBService{
		GetItemFunc: func(ctx context.Context, id string) (*db.Item, error) {
			return &db.Item{ID: id, Name: "testName"}, nil
		},
	}

	response, err := HandleRequestWithService(ctx, request, mockService)

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	// Error case
	mockService = &mockDynamoDBService{
		GetItemFunc: func(ctx context.Context, id string) (*db.Item, error) {
			return nil, errors.New("error getting item")
		},
	}

	_, err = HandleRequestWithService(ctx, request, mockService)

	assert.Error(t, err)
}
