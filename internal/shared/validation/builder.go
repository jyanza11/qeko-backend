package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func BuildErrors(err error) []FieldError {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []FieldError{{
			Field:   "body",
			Code:    "invalid",
			Message: "Invalid request body",
		}}
	}

	out := make([]FieldError, 0, len(validationErrors))
	for _, e := range validationErrors {
		out = append(out, FieldError{
			Field:   e.Field(),
			Code:    e.Tag(),
			Message: Message(e),
		})
	}

	return out
}
