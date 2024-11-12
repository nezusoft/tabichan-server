package trip

import (
	"time"
)

type Trip struct {
	CreatedBy string    `json:"createdBy"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Title     string    `json:"title"`
	ID        string    `json:"id"`
	Completed bool      `json:"completed"`
	Draft     bool      `json:"draft"`
	PlanID    string    `json:"planId"`
}

type Plan struct {
	TripID string
	PlanID string
}

type Itinerary struct {
	ItineraryId   string    `json:"itineraryId"`
	ItineraryName string    `json:"itineraryName"`
	PlanID        string    `json:"planId"`
	TripID        string    `json:"tripId"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
}

type ItineraryItem struct {
	TripID      string
	ItineraryID string
	ID          string
	StartDate   time.Time
	EndDate     time.Time
	Title       string
	Description string
}

type PlanItem struct {
	TripID      string
	PlanID      string
	ID          string
	StartDate   time.Time
	EndDate     time.Time
	Title       string
	Description string
}

type TimeRange struct {
	StartDate time.Time
	EndDate   time.Time
}
