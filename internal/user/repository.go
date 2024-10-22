package user

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tabichanorg/tabichan-server/internal/utils"
)

type UserRepository struct {
	Client *dynamodb.Client
}

func (r *UserRepository) CreateUser(user UserLogin) error {
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

func (r *UserRepository) GetUserByUsernameOrEmail(usernameOrEmailInput string) (*UserLogin, error) {
	if utils.IsEmail(usernameOrEmailInput) {
		return r.GetUserByEmail(usernameOrEmailInput)
	}
	return r.GetUserByUsername(usernameOrEmailInput)
}

func (r *UserRepository) GetUserByUsername(username string) (*UserLogin, error) {
	return r.FetchUserInfo(username, "UsernameIndex", "Username = :username", ":username")
}

func (r *UserRepository) GetUserByEmail(email string) (*UserLogin, error) {
	return r.FetchUserInfo(email, "EmailIndex", "Email = :email", ":email")
}

func (r *UserRepository) FetchUserInfo(usernameOrEmailInput, indexName, keyConditionExpression, keyAttribute string) (*UserLogin, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("Users"),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String(keyConditionExpression),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			keyAttribute: &types.AttributeValueMemberS{Value: usernameOrEmailInput},
		},
	}
	result, err := r.Client.Query(context.TODO(), queryInput)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user UserLogin
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserDetailsByID(id string) (*User, error) {
	return r.FetchUserDetails(id, "UserIDIndex", "UserID = :userid", ":userid")
}

func (r *UserRepository) FetchUserDetails(id, indexName, keyConditionExpression, keyAttribute string) (*User, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("Users"),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String(keyConditionExpression),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			keyAttribute: &types.AttributeValueMemberS{Value: id},
		},
	}
	result, err := r.Client.Query(context.TODO(), queryInput)
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

func (r *UserRepository) CreateSession(session *utils.Session) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Sessions"),
		Item: map[string]types.AttributeValue{
			"SessionID": &types.AttributeValueMemberS{Value: session.SessionID},
			"UserID":    &types.AttributeValueMemberS{Value: session.UserID},
			"Device":    &types.AttributeValueMemberS{Value: session.Device},
			"CreatedAt": &types.AttributeValueMemberS{Value: session.CreatedAt.String()},
			"ExpiresAt": &types.AttributeValueMemberS{Value: session.ExpiresAt.String()},
		},
	}

	_, err := r.Client.PutItem(context.TODO(), input)
	return err
}
