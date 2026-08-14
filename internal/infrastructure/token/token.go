package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/fark-tee/fark-tee-backend/internal/config"
)

// ErrInvalidToken is returned when a token fails signature verification, is
// expired, or was signed with an unexpected algorithm.
var ErrInvalidToken = errors.New("token: invalid token")

// Claims is the subset of a verified access token that callers need.
type Claims struct {
	UserID string
}

type Manager struct {
	cfg *config.Config
}

// @WireSet("Infrastructure")
func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

// Issue creates a signed access token for userID, valid for the configured
// JWT TTL.
func (m *Manager) Issue(userID string) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.cfg.JWT.TTL)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.cfg.JWT.Secret))
}

// Verify parses and validates a signed access token, returning its claims.
func (m *Manager) Verify(tokenString string) (Claims, error) {
	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (any, error) {
		return []byte(m.cfg.JWT.Secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	return Claims{UserID: claims.Subject}, nil
}
