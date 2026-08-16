package dto

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type StoryResponse struct {
	ID        string    `json:"id"`
	PartyID   string    `json:"partyId"`
	UserID    string    `json:"userId"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
}

type StoriesResponse struct {
	Stories []StoryResponse `json:"stories"`
}

type CreateStoryFormData struct {
	Image huma.FormFile `form:"image" contentType:"image/*" required:"true"`
}

type CreateStoryRequest struct {
	PartyID string `path:"partyId" required:"true"`
	RawBody huma.MultipartFormFiles[CreateStoryFormData]
}

type ListStoriesByUserRequest struct {
	PartyID string `path:"partyId" required:"true"`
	UserID  string `path:"userId" required:"true"`
}

type DeleteStoryRequest struct {
	PartyID string `path:"partyId" required:"true"`
	StoryID string `path:"storyId" required:"true"`
}

type DeleteStoryResponse struct{}
