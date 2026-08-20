package entity

import "time"

type TripDirection string

const (
	TripDirectionDepart TripDirection = "DEPART"
	TripDirectionReturn TripDirection = "RETURN"
)

type Trip struct {
	ID          string
	PartyID     string
	UserID      string
	Direction   TripDirection
	Destination Destination
	StartedAt   time.Time

	// Polyline is the OSRM-computed road route from the trip's starting
	// position to Destination, encoded as a Google polyline (precision 5).
	// It is computed once, when the trip starts, rather than recomputed on
	// every position update.
	Polyline string
}
