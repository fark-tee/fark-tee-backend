package handler

import (
	"github.com/fark-tee/fark-tee-backend/internal/handler/instagramoauth"
	"github.com/fark-tee/fark-tee-backend/internal/handler/savedlocation"
)

type Handlers struct {
	InstagramOAuth instagramoauth.Handler
	SavedLocation  savedlocation.Handler
}

// @WireSet("Handler")
func NewHandlers(instagramOAuthHandler instagramoauth.Handler, savedLocationHandler savedlocation.Handler) *Handlers {
	return &Handlers{
		InstagramOAuth: instagramOAuthHandler,
		SavedLocation:  savedLocationHandler,
	}
}
