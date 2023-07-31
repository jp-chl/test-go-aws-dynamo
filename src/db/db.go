package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBOperations interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type AttributeMarshaler interface {
	Marshal(interface{}) (types.AttributeValue, error)
}

type AttributeUnmarshaler interface {
	UnmarshalMap(map[string]types.AttributeValue, interface{}) error
}

type DynamoDBService struct {
	client      DynamoDBOperations
	marshaler   AttributeMarshaler
	unmarshaler AttributeUnmarshaler
}

func NewDynamoDBService(client DynamoDBOperations, marshaler AttributeMarshaler, unmarshaler AttributeUnmarshaler) *DynamoDBService {
	return &DynamoDBService{
		client:      client,
		marshaler:   marshaler,
		unmarshaler: unmarshaler,
	}
}

type Item struct {
	ID   string
	Name string
}

func (s *DynamoDBService) GetItem(ctx context.Context, id string) (*Item, error) {
	av, err := s.marshaler.Marshal(id)
	if err != nil {
		return nil, err
	}

	params := &dynamodb.GetItemInput{
		TableName: aws.String("MyTable"),
		Key: map[string]types.AttributeValue{
			"ID": av,
		},
	}

	resp, err := s.client.GetItem(ctx, params)
	if err != nil {
		return nil, err
	}

	var item Item
	err = s.unmarshaler.UnmarshalMap(resp.Item, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
