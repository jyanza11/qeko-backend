package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appErrors "github.com/jyanza11/qeko-backend/internal/shared/errors"
	"github.com/jyanza11/qeko-backend/internal/shared/response"
	"github.com/jyanza11/qeko-backend/internal/shared/validation"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if !bindAndValidate(c, &req) {
		return
	}

	result, appErr := h.service.Register(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.Success(c, http.StatusCreated, "registered successfully", result)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if !bindAndValidate(c, &req) {
		return
	}

	result, appErr := h.service.Login(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.Success(c, http.StatusOK, "logged in successfully", result)
}

func (h *Handler) Me(c *gin.Context) {
	userID, appErr := UserIDFromContext(c)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	result, appErr := h.service.Me(c.Request.Context(), userID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}

	response.Success(c, http.StatusOK, "ok", result)
}

func bindAndValidate[T any](c *gin.Context, req *T) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		response.Error(c, appErrors.Validation([]validation.FieldError{{
			Field:   "body",
			Code:    "invalid",
			Message: "Invalid request body",
		}}))
		return false
	}

	if err := validation.Validate.Struct(req); err != nil {
		response.Error(c, appErrors.Validation(validation.BuildErrors(err)))
		return false
	}

	return true
}

func UserIDFromContext(c *gin.Context) (uuid.UUID, *appErrors.AppError) {
	value, ok := c.Get(ContextUserIDKey)
	if !ok {
		return uuid.Nil, appErrors.Unauthorized("missing authentication")
	}
	id, ok := value.(uuid.UUID)
	if !ok {
		return uuid.Nil, appErrors.Unauthorized("invalid authentication")
	}
	return id, nil
}
