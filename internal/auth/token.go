package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID         uuid.UUID `json:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Email          string    `json:"email"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret  []byte
	expires time.Duration
}

func NewTokenManager(secret string, expires time.Duration) *TokenManager {
	return &TokenManager{
		secret:  []byte(secret),
		expires: expires,
	}
}

func (m *TokenManager) Generate(userID, organizationID uuid.UUID, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.expires)

	claims := Claims{
		UserID:         userID,
		OrganizationID: organizationID,
		Email:          email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func (m *TokenManager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func (m *TokenManager) ExpiresInSeconds() int64 {
	return int64(m.expires.Seconds())
}
