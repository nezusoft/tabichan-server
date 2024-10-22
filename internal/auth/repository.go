package auth

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

type AuthRepository struct {
	db *dynamodb.Client
}
