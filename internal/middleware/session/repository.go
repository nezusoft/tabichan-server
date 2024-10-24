package middleware

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/tabichanorg/tabichan-server/internal/utils"
)

type MiddlewareRepository struct {
	Client *dynamodb.Client
}

func (r *MiddlewareRepository) GetSession(sessionID string) (*utils.Session, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("Sessions"),
		IndexName:              aws.String("SessionIDIndex"),
		KeyConditionExpression: aws.String("SessionID = :sessionID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sessionID": &types.AttributeValueMemberS{Value: sessionID},
		},
	}
	result, err := r.Client.Query(context.TODO(), queryInput)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("session not found")
	}

	var session utils.Session
	err = attributevalue.UnmarshalMap(result.Items[0], &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *MiddlewareRepository) UpdateSession(oldSessionID, newSessionID, newExpiresAt string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Sessions"),
		Key: map[string]types.AttributeValue{
			"SessionID": &types.AttributeValueMemberS{Value: oldSessionID},
		},
		UpdateExpression: aws.String("SET SessionID = :newSessionID, ExpiresAt = :newExpiry"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newSessionID": &types.AttributeValueMemberS{Value: newSessionID},
			":newExpiry":    &types.AttributeValueMemberS{Value: newExpiresAt},
		},
	}

	_, err := r.Client.UpdateItem(context.TODO(), input)
	return err
}
