package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func init() {
	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")

		if name == "-" {
			return ""
		}

		return name
	})
}
