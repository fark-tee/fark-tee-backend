package handler

import (
	"github.com/fark-tee/fark-tee-backend/internal/handler/instagramoauth"
)

type Handlers struct {
	InstagramOAuth instagramoauth.Handler
}

// @WireSet("Handler")
func NewHandlers(instagramOAuthHandler instagramoauth.Handler) *Handlers {
	return &Handlers{
		InstagramOAuth: instagramOAuthHandler,
	}
}
