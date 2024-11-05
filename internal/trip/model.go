package trip

import (
	"time"
)

type Trip struct {
	UserID    string
	StartDate time.Time
	EndDate   time.Time
	Title     string
	ID        string
	Completed bool
	Draft     bool
	PlanID    string
}

type Plan struct {
	TripId string
}

type Itinerary struct {
	ItineraryName string
	PlanID        string
	TripID        string
	StartDate     time.Time
	EndDate       time.Time
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
