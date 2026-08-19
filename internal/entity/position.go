package entity

import "time"

type Position struct {
	ID         string
	TripID     string
	PartyID    string
	UserID     string
	Lat        float64
	Lng        float64
	RecordedAt time.Time

	// EstimatedDurationSeconds and EstimatedArrivalAt are the OSRM-computed
	// travel time from this position to the trip's destination, as of the
	// moment this position was recorded.
	EstimatedDurationSeconds int
	EstimatedArrivalAt       time.Time
}
