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
}
