package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jyanza11/qeko-backend/internal/auth"
	"github.com/jyanza11/qeko-backend/internal/health"
)

// Handlers groups HTTP handlers. Add new modules as fields when the API grows.
type Handlers struct {
	Health *health.Handler
	Auth   *auth.Handler
}

type Middleware struct {
	Authenticate gin.HandlerFunc
}

func NewRouter(h Handlers, mw Middleware) *gin.Engine {
	r := gin.Default()

	r.GET("/health", h.Health.Check)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", h.Auth.Register)
		authGroup.POST("/login", h.Auth.Login)
		authGroup.GET("/me", mw.Authenticate, h.Auth.Me)
	}

	return r
}
