package entity

type User struct {
	ID              string
	ProfileImageURL string
	DisplayName     string
	Username        string
	GoogleUserID    string
	Rating          float64
	RatingCount     int
	OnTimeCount     int
	LateCount       int

	// EmergencyContactName/Phone are shown, with a tap-to-call action, to
	// every other accepted party member when this user answers "not okay"
	// to a check-in request. Optional - a user who hasn't set one yet just
	// shows up with no contact info in that alert.
	EmergencyContactName  string
	EmergencyContactPhone string
}
