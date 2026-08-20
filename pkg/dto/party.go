package dto

import "time"

type PartyResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	DestinationName string    `json:"destinationName"`
	DestinationLat  float64   `json:"destinationLat"`
	DestinationLng  float64   `json:"destinationLng"`
	TargetTime      time.Time `json:"targetTime"`
	CreatedByID     string    `json:"createdById"`
	CreatedByName   string    `json:"createdByName"`
	Note            string    `json:"note,omitempty"`
}

type CreatePartyRequest struct {
	Body struct {
		Name            string    `json:"name" required:"true"`
		DestinationName string    `json:"destinationName" required:"true"`
		DestinationLat  float64   `json:"destinationLat" required:"true"`
		DestinationLng  float64   `json:"destinationLng" required:"true"`
		TargetTime      time.Time `json:"targetTime" required:"true"`
		// Note is an optional free-text reminder from the creator (e.g.
		// "bring cash"), shown to every party member.
		Note string `json:"note,omitempty"`
	}
}

type PartyMemberResponse struct {
	ID               string `json:"id"`
	PartyID          string `json:"partyId"`
	UserID           string `json:"userId"`
	UserDisplayName  string `json:"userDisplayName"`
	UserProfileImage string `json:"userProfileImage"`
	Status           string `json:"status"`
	TripStatus       string `json:"tripStatus"`

	CheckInStatus            string `json:"checkInStatus"`
	CheckInRequestedByUserID string `json:"checkInRequestedByUserId,omitempty"`
}

type InviteToPartyRequest struct {
	PartyID string `path:"partyId" required:"true"`
	Body    struct {
		UserID string `json:"userId" required:"true"`
	}
}

type PartyInviteResponse struct {
	Party  PartyResponse       `json:"party"`
	Member PartyMemberResponse `json:"member"`
}

type PartyInvitesResponse struct {
	Invites []PartyInviteResponse `json:"invites"`
}

type MyInvitesRequest struct{}

type PartiesResponse struct {
	Parties []PartyResponse `json:"parties"`
}

type MyPartiesRequest struct{}

type GetPartyRequest struct {
	PartyID string `path:"partyId" required:"true"`
}

type PartyMembersResponse struct {
	Members []PartyMemberResponse `json:"members"`
}

type ListPartyMembersRequest struct {
	PartyID string `path:"partyId" required:"true"`
}

type AcceptInviteRequest struct {
	PartyID string `path:"partyId" required:"true"`
}

type DeclineInviteRequest struct {
	PartyID string `path:"partyId" required:"true"`
}

type DeclineInviteResponse struct{}

type RemovePartyMemberRequest struct {
	PartyID string `path:"partyId" required:"true"`
	UserID  string `path:"userId" required:"true"`
}

type RemovePartyMemberResponse struct{}

type NudgePartyMemberRequest struct {
	PartyID string `path:"partyId" required:"true"`
	UserID  string `path:"userId" required:"true"`
}

type NudgePartyMemberResponse struct{}

type RequestCheckInRequest struct {
	PartyID string `path:"partyId" required:"true"`
	UserID  string `path:"userId" required:"true"`
}

type RequestCheckInResponse struct{}

type RespondCheckInRequest struct {
	PartyID string `path:"partyId" required:"true"`
	Body    struct {
		Status string `json:"status" enum:"OK,NOT_OK" required:"true"`
	}
}
