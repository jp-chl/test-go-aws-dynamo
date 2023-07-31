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

type mockAttributeMarshaler struct {
	MarshalFunc func(interface{}) (types.AttributeValue, error)
}

func (m *mockAttributeMarshaler) Marshal(in interface{}) (types.AttributeValue, error) {
	return m.MarshalFunc(in)
}

type mockAttributeUnmarshaler struct {
	UnmarshalMapFunc func(map[string]types.AttributeValue, interface{}) error
}

func (m *mockAttributeUnmarshaler) UnmarshalMap(marshaledMap map[string]types.AttributeValue, out interface{}) error {
	return m.UnmarshalMapFunc(marshaledMap, out)
}

func TestGetItem(t *testing.T) {
	ctx := context.TODO()
	id := "123"

	// Mock the DynamoDB client
	m := &mockDynamoDBClient{
		GetItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{
					"ID":   &types.AttributeValueMemberS{Value: id},
					"Name": &types.AttributeValueMemberS{Value: "ItemName"},
				},
			}, nil
		},
	}

	// Mock the Marshaler
	marshaler := &mockAttributeMarshaler{
		MarshalFunc: func(in interface{}) (types.AttributeValue, error) {
			return &types.AttributeValueMemberS{Value: id}, nil
		},
	}

	// Mock the Unmarshaler
	unmarshaler := &mockAttributeUnmarshaler{
		UnmarshalMapFunc: func(marshaledMap map[string]types.AttributeValue, out interface{}) error {
			item := out.(*Item)
			item.ID = id
			item.Name = "ItemName"
			return nil
		},
	}

	service := NewDynamoDBService(m, marshaler, unmarshaler)

	item, err := service.GetItem(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, item.ID)
	assert.Equal(t, "ItemName", item.Name)

	// Test case for Marshaler error
	marshalerWithError := &mockAttributeMarshaler{
		MarshalFunc: func(in interface{}) (types.AttributeValue, error) {
			return nil, errors.New("marshal error")
		},
	}

	service = NewDynamoDBService(m, marshalerWithError, unmarshaler)

	item, err = service.GetItem(ctx, id)
	assert.Nil(t, item)
	assert.Error(t, err)

	// Test case for Unmarshaler error
	unmarshalerWithError := &mockAttributeUnmarshaler{
		UnmarshalMapFunc: func(marshaledMap map[string]types.AttributeValue, out interface{}) error {
			return errors.New("unmarshal error")
		},
	}

	service = NewDynamoDBService(m, marshaler, unmarshalerWithError)

	item, err = service.GetItem(ctx, id)
	assert.Nil(t, item)
	assert.Error(t, err)
}
