package user

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserRepository struct {
	Client *dynamodb.Client
}

func (r *UserRepository) CreateUser(user User) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item: map[string]types.AttributeValue{
			"Username":      &types.AttributeValueMemberS{Value: user.Username},
			"Password":      &types.AttributeValueMemberS{Value: user.Password},
			"Email":         &types.AttributeValueMemberS{Value: user.Email},
			"OAuthProvider": &types.AttributeValueMemberS{Value: user.OAuthProvider},
			"UserID":        &types.AttributeValueMemberS{Value: user.ID},
		},
	}

	_, err := r.Client.PutItem(context.TODO(), input)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (*User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Users"),
		IndexName:              aws.String("UsernameIndex"),
		KeyConditionExpression: aws.String("Username = :username"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username": &types.AttributeValueMemberS{Value: username},
		},
	}

	result, err := r.Client.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
