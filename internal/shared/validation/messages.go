package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Message(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required", err.Field())
	case "email":
		return fmt.Sprintf("The %s field must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("The %s field must be at least %s characters long", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("The %s field must be at most %s characters long", err.Field(), err.Param())
	case "unique":
		return fmt.Sprintf("The %s field must be unique", err.Field())
	case "exists":
		return fmt.Sprintf("The %s field must exist", err.Field())
	case "confirmed":
		return fmt.Sprintf("The %s field confirmation does not match", err.Field())
	case "alpha":
		return fmt.Sprintf("The %s field must contain only alphabetic characters", err.Field())
	default:
		return err.Error()
	}
}
