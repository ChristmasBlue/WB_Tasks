package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

func CreateValidationErrorResponse(err error) string {
	var msg string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			switch fe.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", fe.Field())
			default:
				msg = fmt.Sprintf("%s is not valid due to %s", fe.Field(), fe.Tag())
			}
		}
	}
	return msg
}
