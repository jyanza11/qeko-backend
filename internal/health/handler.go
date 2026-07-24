package health

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jyanza11/qeko-backend/internal/shared/errors"
	"github.com/jyanza11/qeko-backend/internal/shared/response"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	db    Pinger
	redis Pinger
}

func NewHandler(db, redis Pinger) *Handler {
	return &Handler{db: db, redis: redis}
}

func (h *Handler) Check(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.db.Ping(ctx); err != nil {
		response.Error(c, errors.Internal("database is not running"))
		return
	}

	if err := h.redis.Ping(ctx); err != nil {
		response.Error(c, errors.Internal("redis is not running"))
		return
	}

	response.Success(c, http.StatusOK, "qeko api is running successfully", map[string]string{
		"database": "ok",
		"redis":    "ok",
	})
}
