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
}
