package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DynamoDBOperations interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type DynamoDBService struct {
	Client DynamoDBOperations
}

func NewDynamoDBService(client DynamoDBOperations) *DynamoDBService {
	return &DynamoDBService{
		Client: client,
	}
}

func (d *DynamoDBService) GetItem(ctx context.Context, id string) (*Item, error) {
	av, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String("MyTable"),
		Key: map[string]types.AttributeValue{
			"ID": av,
		},
	}

	resp, err := d.Client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	item := Item{}
	err = attributevalue.UnmarshalMap(resp.Item, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
