package trip

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/tabichanorg/tabichan-server/internal/utils"
)

type TripRepository struct {
	Client *dynamodb.Client
}

func (r *TripRepository) GetTrips(userID string) ([]*Trip, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("Trips"),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
		},
	}
	result, err := r.Client.Query(context.TODO(), queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trips for user with ID %s: %w", userID, err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("couldn't find any trips for user with ID: %s", userID)
	}

	var trips []*Trip
	for _, item := range result.Items {
		var trip Trip
		err = attributevalue.UnmarshalMap(item, &trip)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal trip: %w", err)
		}
		trips = append(trips, &trip)
	}

	return trips, nil
}

func (r *TripRepository) GetTrip(tripID string) (*Trip, error) {
	queryInput := &dynamodb.GetItemInput{
		TableName: aws.String("Trips"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("TRIP#%s", tripID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("META#%s", tripID)},
		},
	}
	result, err := r.Client.GetItem(context.TODO(), queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trip with ID %s: %w", tripID, err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("trip doesn't exist")
	}

	var trip Trip
	err = attributevalue.UnmarshalMap(result.Item, &trip)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal trip: %w", err)
	}

	return &trip, nil
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
			"StartDate": &types.AttributeValueMemberS{Value: formatTime(tripData.StartDate)},
			"EndDate":   &types.AttributeValueMemberS{Value: formatTime(tripData.EndDate)},
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

// TODO:
func (r *TripRepository) EditTrip(tripID string, tripData *Trip) error {
	// Define the update expression and attribute values
	updateExpression := "SET "
	attributeValues := map[string]types.AttributeValue{}
	attributeNames := map[string]string{}

	// Only add fields that are set in the tripData
	if !tripData.StartDate.IsZero() {
		updateExpression += "#StartDate = :startDate, "
		attributeValues[":startDate"] = &types.AttributeValueMemberS{Value: tripData.StartDate.Format(time.RFC3339)}
		attributeNames["#StartDate"] = "StartDate"
	}
	if !tripData.EndDate.IsZero() {
		updateExpression += "#EndDate = :endDate, "
		attributeValues[":endDate"] = &types.AttributeValueMemberS{Value: tripData.EndDate.Format(time.RFC3339)}
		attributeNames["#EndDate"] = "EndDate"
	}
	if tripData.Title != "" {
		updateExpression += "#Title = :title, "
		attributeValues[":title"] = &types.AttributeValueMemberS{Value: tripData.Title}
		attributeNames["#Title"] = "Title"
	}
	if tripData.Completed {
		updateExpression += "#Completed = :completed, "
		attributeValues[":completed"] = &types.AttributeValueMemberBOOL{Value: tripData.Completed}
		attributeNames["#Completed"] = "Completed"
	}
	if tripData.Draft {
		updateExpression += "#Draft = :draft, "
		attributeValues[":draft"] = &types.AttributeValueMemberBOOL{Value: tripData.Draft}
		attributeNames["#Draft"] = "Draft"
	}

	// If there are no fields to update, return an error
	if len(attributeValues) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Remove the trailing comma and space from the update expression
	updateExpression = updateExpression[:len(updateExpression)-2]

	// Prepare the update input for DynamoDB
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Trips"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("TRIP#%s", tripID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("META#%s", tripID)},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: attributeValues,
		ExpressionAttributeNames:  attributeNames,
	}

	_, err := r.Client.UpdateItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	return nil
}

func (r *TripRepository) DeleteTrip(tripID string) error {
	queryInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("Trips"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("TRIP#%s", tripID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("META#%s", tripID)},
		},
	}
	_, err := r.Client.DeleteItem(context.TODO(), queryInput)
	if err != nil {
		return fmt.Errorf("failed to delete trip with ID %s: %w", tripID, err)
	}

	return nil
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

func (r *TripRepository) GetItinerary(itineraryID string) (*Itinerary, error) {
	queryInput := &dynamodb.GetItemInput{
		TableName: aws.String("Itineraries"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ITINERARY#%s", itineraryID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("META#%s", itineraryID)},
		},
	}
	result, err := r.Client.GetItem(context.TODO(), queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch itinerary with ID %s: %w", itineraryID, err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("itinerary doesn't exist")
	}

	var itinerary Itinerary
	err = attributevalue.UnmarshalMap(result.Item, &itinerary)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal itinerary: %w", err)
	}

	return &itinerary, nil
}
func (r *TripRepository) DeleteItinerary(itineraryID string) error {
	queryInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("Itineraries"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ITINERARY#%s", itineraryID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("META#%s", itineraryID)},
		},
	}
	_, err := r.Client.DeleteItem(context.TODO(), queryInput)
	if err != nil {
		return fmt.Errorf("failed to delete itinerary with ID %s: %w", itineraryID, err)
	}

	return nil
}

func (r *TripRepository) CreateItinerary(createItineraryData Itinerary) (*Itinerary, error) {
	createItineraryData.ItineraryId = utils.GenerateID()
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Itineraries"),
		Item: map[string]types.AttributeValue{
			"PK":            &types.AttributeValueMemberS{Value: fmt.Sprintf("ITINERARY#%s", createItineraryData.ItineraryId)},
			"SK":            &types.AttributeValueMemberS{Value: fmt.Sprintf("META#%s", createItineraryData.ItineraryId)},
			"GSI1PK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("PLAN#%s", createItineraryData.PlanID)},
			"GSI2PK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("TRIP#%s", createItineraryData.TripID)},
			"PlanID":        &types.AttributeValueMemberS{Value: createItineraryData.PlanID},
			"TripID":        &types.AttributeValueMemberS{Value: createItineraryData.TripID},
			"StartDate":     &types.AttributeValueMemberS{Value: formatTime(createItineraryData.StartDate)},
			"EndDate":       &types.AttributeValueMemberS{Value: formatTime(createItineraryData.EndDate)},
			"ItineraryName": &types.AttributeValueMemberS{Value: createItineraryData.ItineraryName},
			"ItineraryId":   &types.AttributeValueMemberS{Value: createItineraryData.ItineraryId},
		},
	}

	_, err := r.Client.PutItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return &createItineraryData, err
}

func formatTime(date time.Time) string {
	return date.Format(time.RFC3339)
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
