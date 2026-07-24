package response

import (
	"github.com/gin-gonic/gin"
	"github.com/jyanza11/qeko-backend/internal/shared/errors"
)

func Error(c *gin.Context, err *errors.AppError) {

	c.JSON(err.StatusCode, Response[any]{
		Status:  "error",
		Message: err.Message,
		Code:    err.Code,
		Errors:  err.Details,
		Data:    nil,
	})
}
