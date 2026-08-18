package dto

import "time"

type TripResponse struct {
	ID        string    `json:"id"`
	PartyID   string    `json:"partyId"`
	UserID    string    `json:"userId"`
	Direction string    `json:"direction"`
	StartedAt time.Time `json:"startedAt"`
}

type PositionResponse struct {
	ID         string    `json:"id"`
	TripID     string    `json:"tripId"`
	PartyID    string    `json:"partyId"`
	UserID     string    `json:"userId"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	RecordedAt time.Time `json:"recordedAt"`
}

type PositionsResponse struct {
	Positions []PositionResponse `json:"positions"`
}

type StartTripResponse struct {
	Trip     TripResponse     `json:"trip"`
	Position PositionResponse `json:"position"`
}

type StartTripRequest struct {
	PartyID string `path:"partyId" required:"true"`
	Body    struct {
		Direction string  `json:"direction" enum:"DEPART,RETURN" required:"true"`
		Lat       float64 `json:"lat" required:"true"`
		Lng       float64 `json:"lng" required:"true"`
	}
}

type UpdatePositionRequest struct {
	PartyID string `path:"partyId" required:"true"`
	Body    struct {
		Lat float64 `json:"lat" required:"true"`
		Lng float64 `json:"lng" required:"true"`
	}
}

type GetMemberPositionRequest struct {
	PartyID string `path:"partyId" required:"true"`
	UserID  string `path:"userId" required:"true"`
}

type GetPartyPositionsRequest struct {
	PartyID string `path:"partyId" required:"true"`
}

type UpdateTripStatusRequest struct {
	PartyID string `path:"partyId" required:"true"`
	Body    struct {
		Status string `json:"status" enum:"PENDING_DEPARTURE,DEPARTED,ARRIVED,RETURNING,RETURNED" required:"true"`
	}
}
