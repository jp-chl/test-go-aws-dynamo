package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetItem(ctx context.Context, id string) (*Item, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)

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

	resp, err := client.GetItem(ctx, input)
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
