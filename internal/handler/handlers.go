package handler

import (
	"github.com/fark-tee/fark-tee-backend/internal/handler/devicetoken"
	"github.com/fark-tee/fark-tee-backend/internal/handler/googleoauth"
	"github.com/fark-tee/fark-tee-backend/internal/handler/party"
	"github.com/fark-tee/fark-tee-backend/internal/handler/review"
	"github.com/fark-tee/fark-tee-backend/internal/handler/savedlocation"
	"github.com/fark-tee/fark-tee-backend/internal/handler/story"
	"github.com/fark-tee/fark-tee-backend/internal/handler/trip"
	"github.com/fark-tee/fark-tee-backend/internal/handler/upload"
	"github.com/fark-tee/fark-tee-backend/internal/handler/user"
)

type Handlers struct {
	GoogleOAuth   googleoauth.Handler
	SavedLocation savedlocation.Handler
	Party         party.Handler
	Story         story.Handler
	Trip          trip.Handler
	Upload        upload.Handler
	User          user.Handler
	DeviceToken   devicetoken.Handler
	Review        review.Handler
}

// @WireSet("Handler")
func NewHandlers(
	googleOAuthHandler googleoauth.Handler,
	savedLocationHandler savedlocation.Handler,
	partyHandler party.Handler,
	storyHandler story.Handler,
	tripHandler trip.Handler,
	uploadHandler upload.Handler,
	userHandler user.Handler,
	deviceTokenHandler devicetoken.Handler,
	reviewHandler review.Handler,
) *Handlers {
	return &Handlers{
		GoogleOAuth:   googleOAuthHandler,
		SavedLocation: savedLocationHandler,
		Party:         partyHandler,
		Story:         storyHandler,
		Trip:          tripHandler,
		Upload:        uploadHandler,
		User:          userHandler,
		DeviceToken:   deviceTokenHandler,
		Review:        reviewHandler,
	}
}
