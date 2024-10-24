package db

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/tabichanorg/tabichan-server/internal/utils"
)

func GetSession(sessionID string) (*utils.Session, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Sessions"),
		Key: map[string]*dynamodb.AttributeValue{
			"SessionID": {
				S: aws.String(sessionID),
			},
		},
	}

	result, err := svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("session not found")
	}

	session := &utils.Session{
		SessionID: *result.Item["SessionID"].S,
		UserID:    *result.Item["UserID"].S,
		ExpiresAt: *result.Item["ExpiresAt"].S,
	}

	return session, nil
}
