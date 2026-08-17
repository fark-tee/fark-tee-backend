package entity

type User struct {
	ID              string
	ProfileImageURL string
	DisplayName     string
	GoogleUserID    string
	Rating          float64
	RatingCount     int
	OnTimeCount     int
	LateCount       int
}
