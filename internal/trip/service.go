package trip

import (
	"fmt"
	"time"

	"github.com/tabichanorg/tabichan-server/internal/utils"
)

type TripService struct {
	Repo *TripRepository
}

func (s *TripService) GetTrips(userID string) ([]*Trip, error) {
	trips, err := s.Repo.GetTrips(userID)
	if err != nil {
		return nil, fmt.Errorf(`error fetching trips for user with id %s: %w`, userID, err)
	}

	return trips, nil
}

func (s *TripService) GetTrip(tripID string) (*Trip, error) {
	trip, err := s.Repo.GetTrip(tripID)
	if err != nil {
		return nil, fmt.Errorf(`error fetching trip with id %s: %s`, tripID, err)
	}

	return trip, nil
}

func (s *TripService) DeleteTrip(tripID string) error {
	err := s.Repo.DeleteTrip(tripID)
	if err != nil {
		return fmt.Errorf(`error deleting trip with id %s: %s`, tripID, err)
	}

	return nil
}

func (s *TripService) CreateTrip(tripData *Trip) (*Trip, error) {

	if err := validateTripData(tripData); err != nil {
		return nil, fmt.Errorf(`error validating trip data: %s`, err)
	}

	tripID := utils.GenerateID()
	planID := utils.GenerateID()

	plan, err := s.CreatePlan(planID, tripID)
	if err != nil {
		return nil, fmt.Errorf(`error creating plan: "%s"`, err)
	}

	if err := s.Repo.CreateTrip(tripData, plan); err != nil {
		return nil, fmt.Errorf(`error creating trip: "%s"`, err)
	}

	tripData.ID = tripID
	tripData.PlanID = planID

	return tripData, err
}

func (s *TripService) CreatePlan(planID, tripID string) (*Plan, error) {
	plan, err := s.Repo.CreatePlan(planID, tripID)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *TripService) GetItineraries(planId string) ([]*Itinerary, error) {
	itineraries, err := s.Repo.GetItineraries(planId)
	if err != nil {
		return nil, fmt.Errorf(`error fetching itineraries for trip with plan id %s: %w`, planId, err)
	}

	return itineraries, nil
}

func (s *TripService) GetItinerary(itineraryId string) (*Itinerary, error) {
	itinerary, err := s.Repo.GetItinerary(itineraryId)
	if err != nil {
		return nil, fmt.Errorf(`error fetching itinerary with id %s: %s`, itineraryId, err)
	}

	return itinerary, nil
}

func (s *TripService) DeleteItinerary(itineraryId string) error {
	err := s.Repo.DeleteItinerary(itineraryId)
	if err != nil {
		return fmt.Errorf(`error deleting itinerary with id %s: %s`, itineraryId, err)
	}

	return nil
}

func (s *TripService) CreateItinerary(createItineraryData Itinerary) (*Itinerary, error) {
	trip, err := s.GetTrip(createItineraryData.TripID)
	if err != nil {
		return nil, err
	}

	rangeOne := TimeRange{trip.StartDate, trip.EndDate}
	rangeTwo := TimeRange{createItineraryData.StartDate, createItineraryData.EndDate}
	if err := validateDatesWithinRange(rangeOne, rangeTwo); err != nil {
		return nil, err
	}

	createItineraryData.PlanID = trip.PlanID
	itinerary, err := s.Repo.CreateItinerary(createItineraryData)
	if err != nil {
		return nil, err
	}
	return itinerary, nil
}

func (s *TripService) GetItineraryItems(itineraryId string) ([]*ItineraryItem, error) {
	itineraryItems, err := s.Repo.GetItineraryItems(itineraryId)

	if err != nil {
		return nil, fmt.Errorf(`error fetching items for itinerary with itinerary id %s: %w`, itineraryId, err)

	}
	return itineraryItems, nil
}

func (s *TripService) CreateItineraryItem(createItineraryItemData ItineraryItem) (*ItineraryItem, error) {
	// validate createItineraryItemData dates
	if !createItineraryItemData.StartDate.Before(createItineraryItemData.EndDate) {
		return nil, fmt.Errorf("start date before end date")
	}
	// needs to verify itinerary start/end date is still within trip start/end date
	trip, err := s.GetTrip(createItineraryItemData.TripID)
	if err != nil {
		return nil, err
	}
	
	// validate createItineraryItemData within trip dates
	rangeOne := TimeRange{trip.StartDate, trip.EndDate}
	rangeTwo := TimeRange{createItineraryItemData.StartDate, createItineraryItemData.EndDate}
	if err := validateDatesWithinRange(rangeOne, rangeTwo); err != nil {
		return nil, err
	}

	if len(createItineraryItemData.Title) > 15 {
		return nil, fmt.Errorf("title must be a maximum of 15 characters long")
	}

	if len(createItineraryItemData.Description) > 100 {
		return nil, fmt.Errorf("description must be a maximum of 100 characters long")
	}


	// get itinerary
	itinerary, err := s.GetItinerary(createItineraryItemData.ItineraryID)
	if err != nil {
		return nil, err
	}

	// edit itinerary start and end dates
	if createItineraryItemData.StartDate.Before(itinerary.StartDate) {
		// edit endpoint to change itinerary data
	}
	if createItineraryItemData.EndDate.After(itinerary.EndDate) {
		// edit endpoint to change itinerary data
	}

	itineraryItem, err := s.Repo.CreateItineraryItem(createItineraryItemData)
	if err != nil {
		return nil, err
	}
	return itineraryItem, nil
}

func validateTripData(tripData *Trip) error {
	if !tripData.StartDate.Before(tripData.EndDate) {
		return fmt.Errorf("start date must be before end date")
	}

	if tripData.StartDate.Before(time.Now().UTC()) {
		return fmt.Errorf("start date must be in the future")
	}

	if len(tripData.Title) > 15 {
		return fmt.Errorf("title must be a maximum of 15 characters long")
	}

	return nil
}

func validateDatesWithinRange(rangeOne, rangeTwo TimeRange) error {
	if rangeOne.StartDate.After(rangeTwo.StartDate) {
		return fmt.Errorf("start date must come before parent start date")
	}
	if rangeOne.EndDate.Before(rangeTwo.EndDate) {
		return fmt.Errorf("end date must come before parent end date")
	}
	return nil
}
