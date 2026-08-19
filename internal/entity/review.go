package entity

import "time"

// Review is one party member's rating of another after the target has
// arrived at the destination.
type Review struct {
	ID           string
	PartyID      string
	ReviewerID   string
	TargetUserID string
	Score        int
	CreatedAt    time.Time
}
