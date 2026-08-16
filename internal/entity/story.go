package entity

import "time"

type Story struct {
	ID        string
	PartyID   string
	UserID    string
	Image     string
	CreatedAt time.Time
}
