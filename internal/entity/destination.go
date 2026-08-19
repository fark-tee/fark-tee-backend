package entity

// Destination is the named place a trip is heading to - the party's venue
// for a DEPART trip, or wherever the user chose to go for a RETURN trip.
type Destination struct {
	Name string
	Lat  float64
	Lng  float64
}
