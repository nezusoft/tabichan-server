package db

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DynamoClient *dynamodb.Client

func InitDynamoDB() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	DynamoClient = dynamodb.NewFromConfig(cfg)
	log.Println("Dynamo Client Initiated")
}
