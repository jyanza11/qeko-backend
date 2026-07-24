package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/jyanza11/qeko-backend/internal/shared/errors"
	"github.com/jyanza11/qeko-backend/internal/shared/response"
)

const (
	ContextUserIDKey         = "auth_user_id"
	ContextOrganizationIDKey = "auth_organization_id"
	ContextEmailKey          = "auth_email"
)

func Middleware(tokens *TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, apperrors.Unauthorized("missing authorization header"))
			c.Abort()
			return
		}

		scheme, token, ok := strings.Cut(header, " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
			response.Error(c, apperrors.Unauthorized("invalid authorization header"))
			c.Abort()
			return
		}

		claims, err := tokens.Parse(token)
		if err != nil {
			response.Error(c, apperrors.Unauthorized("invalid or expired token"))
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextOrganizationIDKey, claims.OrganizationID)
		c.Set(ContextEmailKey, claims.Email)
		c.Next()
	}
}

func OrganizationIDFromContext(c *gin.Context) (uuid.UUID, error) {
	value, ok := c.Get(ContextOrganizationIDKey)
	if !ok {
		return uuid.Nil, apperrors.Unauthorized("missing authentication")
	}
	id, ok := value.(uuid.UUID)
	if !ok {
		return uuid.Nil, apperrors.Unauthorized("invalid authentication")
	}
	return id, nil
}
