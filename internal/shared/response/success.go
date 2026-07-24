package response

import "github.com/gin-gonic/gin"

func Success[T any](c *gin.Context, statusCode int, message string, data T) {
	c.JSON(statusCode, Response[T]{
		Status:  "success",
		Message: message,
		Errors:  []any{},
		Data:    data,
	})

}
