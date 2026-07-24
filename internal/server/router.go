package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jyanza11/qeko-backend/internal/health"
)

// Handlers groups HTTP handlers. Add new modules as fields when the API grows.
type Handlers struct {
	Health *health.Handler
}

func NewRouter(h Handlers) *gin.Engine {
	r := gin.Default()

	r.GET("/health", h.Health.Check)

	return r
}
