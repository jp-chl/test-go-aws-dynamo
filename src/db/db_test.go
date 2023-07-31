package db

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type mockDynamoDBClient struct {
	GetItemFunc func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

func (m *mockDynamoDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.GetItemFunc(ctx, params, optFns...)
}

func TestGetItem(t *testing.T) {
	ctx := context.TODO()
	testId := "testID"
	testName := "testName"

	// Success case
	mockClient := &mockDynamoDBClient{
		GetItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			item := map[string]types.AttributeValue{
				"ID":   &types.AttributeValueMemberS{Value: testId},
				"Name": &types.AttributeValueMemberS{Value: testName},
			}
			return &dynamodb.GetItemOutput{Item: item}, nil
		},
	}

	service := NewDynamoDBService(mockClient)
	item, err := service.GetItem(ctx, testId)

	assert.NoError(t, err)
	assert.Equal(t, testId, item.ID)
	assert.Equal(t, testName, item.Name)

	// Error case
	mockClient = &mockDynamoDBClient{
		GetItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return nil, errors.New("error getting item")
		},
	}

	service = NewDynamoDBService(mockClient)
	_, err = service.GetItem(ctx, testId)

	assert.Error(t, err)
}
