package entity

import "time"

type Party struct {
	ID              string
	Name            string
	DestinationName string
	DestinationLat  float64
	DestinationLng  float64
	TargetTime      time.Time
	CreatedByID     string
	CreatedByName   string
}
