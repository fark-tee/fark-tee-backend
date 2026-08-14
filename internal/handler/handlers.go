package handler

import (
	"github.com/fark-tee/fark-tee-backend/internal/handler/instagramoauth"
	"github.com/fark-tee/fark-tee-backend/internal/handler/party"
	"github.com/fark-tee/fark-tee-backend/internal/handler/savedlocation"
	"github.com/fark-tee/fark-tee-backend/internal/handler/user"
)

type Handlers struct {
	InstagramOAuth instagramoauth.Handler
	SavedLocation  savedlocation.Handler
	Party          party.Handler
	User           user.Handler
}

// @WireSet("Handler")
func NewHandlers(
	instagramOAuthHandler instagramoauth.Handler,
	savedLocationHandler savedlocation.Handler,
	partyHandler party.Handler,
	userHandler user.Handler,
) *Handlers {
	return &Handlers{
		InstagramOAuth: instagramOAuthHandler,
		SavedLocation:  savedLocationHandler,
		Party:          partyHandler,
		User:           userHandler,
	}
}
