package authmw

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/apperror"

	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/token"
)

type contextKey string

const userIDContextKey contextKey = "userID"

type Middleware struct {
	tokenManager *token.Manager
}

// @WireSet("Infrastructure")
func New(tokenManager *token.Manager) *Middleware {
	return &Middleware{tokenManager: tokenManager}
}

// RequireAuth returns a huma middleware that rejects requests without a
// valid `Authorization: Bearer <token>` header, and otherwise makes the
// token's user ID available to handlers via UserIDFromContext. api is used
// to write a shaped error response when auth fails; pass the group the
// middleware is attached to.
func (m *Middleware) RequireAuth(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		tokenString, ok := strings.CutPrefix(ctx.Header("Authorization"), "Bearer ")
		if !ok || tokenString == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing bearer token",
				apperror.NewUnauthorizedError("UNAUTHORIZED", "missing bearer token"))

			return
		}

		claims, err := m.tokenManager.Verify(tokenString)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid or expired token",
				apperror.NewUnauthorizedError("UNAUTHORIZED", "invalid or expired token", err))

			return
		}

		next(huma.WithValue(ctx, userIDContextKey, claims.UserID))
	}
}

// UserIDFromContext returns the authenticated user's ID, as set by
// RequireAuth. ok is false if no request in flight was authenticated.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)

	return userID, ok
}

// RequireUserID is a handler-side convenience wrapper around
// UserIDFromContext that returns a ready-to-return apperror when the
// context has no authenticated user - which should not happen for any
// handler registered behind RequireAuth, but keeps handlers safe if a route
// is ever wired up outside that group by mistake.
func RequireUserID(ctx context.Context) (string, error) {
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		return "", apperror.NewUnauthorizedError("UNAUTHORIZED", "missing authenticated user")
	}

	return userID, nil
}
