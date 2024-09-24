package app

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
	"github.com/tabichanorg/tabichan-server/internal/db"
	"github.com/tabichanorg/tabichan-server/internal/server"
)

func InitializeApp() (*server.Server, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db.InitDynamoDB()
	verifyDynamoDBConnection(db.DynamoClient)

	srv := server.NewServer("localhost:8080")

	return srv, nil
}

func verifyDynamoDBConnection(client *dynamodb.Client) {
	result, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalf("Failed to connect to DynamoDB: %v", err)
	}
	fmt.Printf("Successfully connected to DynamoDB. Tables: %v\n", result.TableNames)
}
