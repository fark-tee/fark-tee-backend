package entity

type PartyMemberStatus string

const (
	PartyMemberStatusPending  PartyMemberStatus = "PENDING"
	PartyMemberStatusAccepted PartyMemberStatus = "ACCEPTED"
)

type PartyMember struct {
	ID              string
	PartyID         string
	UserID          string
	UserDisplayName string
	Status          PartyMemberStatus
}
