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

// CheckInStatus tracks the "are you okay?" safety check a party member can
// ask of another member who is RETURNING (heading home). Unlike TripStatus
// it isn't forward-only: a fresh request always resets it back to PENDING,
// and PENDING moves to exactly one of OK/NOT_OK when the target responds.
type CheckInStatus string

const (
	CheckInStatusNone    CheckInStatus = "NONE"
	CheckInStatusPending CheckInStatus = "PENDING"
	CheckInStatusOK      CheckInStatus = "OK"
	CheckInStatusNotOK   CheckInStatus = "NOT_OK"
)

type PartyMember struct {
	ID               string
	PartyID          string
	UserID           string
	UserDisplayName  string
	UserProfileImage string
	Status           PartyMemberStatus
	TripStatus       TripStatus

	CheckInStatus            CheckInStatus
	CheckInRequestedByUserID string
}
