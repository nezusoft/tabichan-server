package trip

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type TripRepository struct {
	Client *dynamodb.Client
}

func (r *TripRepository) CreateTrip(tripData *Trip, planData *Plan) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Trips"),
		Item: map[string]types.AttributeValue{
			"PK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("TRIP#%s", planData.TripID)},
			"SK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("META#%s", planData.TripID)},
			"GSI1PK":    &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", tripData.CreatedBy)},
			"GSI2PK":    &types.AttributeValueMemberS{Value: fmt.Sprintf("TRIP#%s", planData.TripID)},
			"CreatedBy": &types.AttributeValueMemberS{Value: tripData.CreatedBy},
			"StartDate": &types.AttributeValueMemberS{Value: tripData.StartDate.Format(time.RFC3339)},
			"EndDate":   &types.AttributeValueMemberS{Value: tripData.EndDate.Format(time.RFC3339)},
			"Title":     &types.AttributeValueMemberS{Value: tripData.Title},
			"ID":        &types.AttributeValueMemberS{Value: planData.TripID},
			"Completed": &types.AttributeValueMemberBOOL{Value: tripData.Completed},
			"Draft":     &types.AttributeValueMemberBOOL{Value: tripData.Draft},
			"PlanID":    &types.AttributeValueMemberS{Value: planData.PlanID},
		},
	}

	_, err := r.Client.PutItem(context.TODO(), input)
	return err
}
func (r *TripRepository) CreatePlan(planID, tripID string) (*Plan, error) {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Plans"),
		Item: map[string]types.AttributeValue{
			"PK":     &types.AttributeValueMemberS{Value: fmt.Sprintf("TRIP#%s", tripID)},
			"SK":     &types.AttributeValueMemberS{Value: fmt.Sprintf("PLAN#%s", planID)},
			"PlanID": &types.AttributeValueMemberS{Value: planID},
			"TripID": &types.AttributeValueMemberS{Value: tripID},
		},
	}

	_, err := r.Client.PutItem(context.TODO(), input)
	if err != nil {
		return &Plan{}, err
	}

	newPlan := &Plan{
		PlanID: planID,
		TripID: tripID,
	}
	return newPlan, err
}

// Tables

// 1. Trips
// {
// 	"PK": "TRIP#ID"
// 	"SK": "META#ID"
// 	"GSI1PK": "USER#ID"
// 	"GSI1SK": "TRIP#ID"
// }

// 2. Users

// 3. Trip <-> Users
// {
// 	"PK": "USER#ID" // get trips by user
// 	"SK": "TRIP#ID"
// 	"GSI1PK": "TRIP#ID" // get user by trips
// 	"GSI1SK": "USER#ID"
// }
