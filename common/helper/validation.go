package helper

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func TranslateValidationError(err error) map[string]string {
	details := make(map[string]string)

	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()

			switch fieldError.Tag() {
			case "required":
				details[field] = fmt.Sprintf("%s is required", field)

			case "email":
				details[field] = "Email format is invalid"

			case "min":
				details[field] = fmt.Sprintf(
					"%s must contain at least %s characters",
					field,
					fieldError.Param(),
				)

			default:
				details[field] = fmt.Sprintf("%s is invalid", field)
			}
		}

		return details
	}

	details["request"] = "Invalid request data"

	return details
}
