package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/fark-tee/fark-tee-backend/internal/config"
)

// ErrInvalidToken is returned when a token fails signature verification, is
// expired, was signed with an unexpected algorithm, or is not of the type
// the caller asked to verify (e.g. a refresh token presented as an access
// token).
var ErrInvalidToken = errors.New("token: invalid token")

// Claims is the subset of a verified token that callers need.
type Claims struct {
	UserID string
}

type tokenType string

const (
	accessTokenType  tokenType = "access"
	refreshTokenType tokenType = "refresh"
)

// customClaims tags every token with its type so a refresh token can never
// be verified as an access token, or vice versa.
type customClaims struct {
	jwt.RegisteredClaims
	Type tokenType `json:"typ"`
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
	return m.issue(userID, accessTokenType, m.cfg.JWT.TTL)
}

// IssueRefreshToken creates a signed, long-lived refresh token for userID
// that can be redeemed for a new access/refresh token pair.
func (m *Manager) IssueRefreshToken(userID string) (string, error) {
	return m.issue(userID, refreshTokenType, m.cfg.JWT.RefreshTTL)
}

func (m *Manager) issue(userID string, typ tokenType, ttl time.Duration) (string, error) {
	now := time.Now()

	claims := customClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
		Type: typ,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.cfg.JWT.Secret))
}

// Verify parses and validates a signed access token, returning its claims.
func (m *Manager) Verify(tokenString string) (Claims, error) {
	return m.verify(tokenString, accessTokenType)
}

// VerifyRefreshToken parses and validates a signed refresh token, returning
// its claims.
func (m *Manager) VerifyRefreshToken(tokenString string) (Claims, error) {
	return m.verify(tokenString, refreshTokenType)
}

func (m *Manager) verify(tokenString string, want tokenType) (Claims, error) {
	claims := &customClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (any, error) {
		return []byte(m.cfg.JWT.Secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || claims.Type != want {
		return Claims{}, ErrInvalidToken
	}

	return Claims{UserID: claims.Subject}, nil
}
