package entity

type PartyMemberStatus string

const (
	PartyMemberStatusPending  PartyMemberStatus = "PENDING"
	PartyMemberStatusAccepted PartyMemberStatus = "ACCEPTED"
)

// TripStatus tracks a party member's progress through a trip. It only ever
// moves forward - see TripStatusOrder.
type TripStatus string

const (
	TripStatusPendingDeparture TripStatus = "PENDING_DEPARTURE"
	TripStatusDeparted         TripStatus = "DEPARTED"
	TripStatusArrived          TripStatus = "ARRIVED"
	TripStatusReturning        TripStatus = "RETURNING"
	TripStatusReturned         TripStatus = "RETURNED"
)

// tripStatusOrder ranks each TripStatus from least to most advanced.
var tripStatusOrder = map[TripStatus]int{
	TripStatusPendingDeparture: 0,
	TripStatusDeparted:         1,
	TripStatusArrived:          2,
	TripStatusReturning:        3,
	TripStatusReturned:         4,
}

// TripStatusOrder returns s's rank (least to most advanced) and whether s is
// a recognized TripStatus.
func TripStatusOrder(s TripStatus) (int, bool) {
	order, ok := tripStatusOrder[s]
	return order, ok
}

type PartyMember struct {
	ID               string
	PartyID          string
	UserID           string
	UserDisplayName  string
	UserProfileImage string
	Status           PartyMemberStatus
	TripStatus       TripStatus
}
