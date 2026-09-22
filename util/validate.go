package util

import "github.com/go-playground/validator/v10"

func GetErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "len":
		return "Invalid length"
	case "numeric":
		return "Invalid number"
	}
	return "Invalid " + err.Tag()
}
