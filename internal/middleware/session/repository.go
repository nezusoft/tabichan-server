package middleware

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

type MiddlewareRepository struct {
	Client *dynamodb.Client
}
