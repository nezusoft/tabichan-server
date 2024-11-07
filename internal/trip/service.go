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
		return []*Trip{}, fmt.Errorf(`error fetching trips for user with id %s: %w`, userID, err)
	}

	return trips, nil
}

func (s *TripService) GetTrip(tripID string) (*Trip, error) {
	trip, err := s.Repo.GetTrip(tripID)
	if err != nil {
		return &Trip{}, fmt.Errorf(`error fetching trip with id %s: %s`, tripID, err)
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
		return &Trip{}, fmt.Errorf(`error validating trip data: %s`, err)
	}

	tripID := utils.GenerateID()
	planID := utils.GenerateID()

	plan, err := s.CreatePlan(planID, tripID)
	if err != nil {
		return &Trip{}, fmt.Errorf(`error creating plan: "%s"`, err)
	}

	if err := s.Repo.CreateTrip(tripData, plan); err != nil {
		return &Trip{}, fmt.Errorf(`error creating trip: "%s"`, err)
	}

	tripData.ID = tripID
	tripData.PlanID = planID

	return tripData, err
}

func (s *TripService) CreatePlan(planID, tripID string) (*Plan, error) {
	plan, err := s.Repo.CreatePlan(planID, tripID)
	if err != nil {
		return &Plan{}, err
	}
	return plan, nil
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
