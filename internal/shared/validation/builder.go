package validation

import (
	"github.com/go-playground/validator/v10"
)

func BuildErrors(err error) []FieldError {

	validationErrors := err.(validator.ValidationErrors)
	errors := make([]FieldError, 0, len(validationErrors))

	for _, e := range validationErrors {

		errors = append(errors, FieldError{
			Field:   e.Field(),
			Code:    e.Tag(),
			Message: Message(e),
		})
	}

	return errors
}
