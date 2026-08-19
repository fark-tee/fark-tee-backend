package dto

import "time"

type ReviewResponse struct {
	ID           string    `json:"id"`
	PartyID      string    `json:"partyId"`
	ReviewerID   string    `json:"reviewerId"`
	TargetUserID string    `json:"targetUserId"`
	Score        int       `json:"score"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"createdAt"`
}

type ReviewsResponse struct {
	Reviews []ReviewResponse `json:"reviews"`
}

type CreateReviewRequest struct {
	PartyID string `path:"partyId" required:"true"`
	UserID  string `path:"userId" required:"true"`
	Body    struct {
		Score   int    `json:"score" required:"true" minimum:"1" maximum:"5"`
		Comment string `json:"comment,omitempty" maxLength:"500"`
	}
}

type ListMyReviewsRequest struct {
	PartyID string `path:"partyId" required:"true"`
}
