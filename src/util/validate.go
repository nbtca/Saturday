package util

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Location int

const (
	BODY Location = iota
	PARAM
	Query
)

func ValidationHandler(structToValidate interface{}) validator.ValidationErrors {
	validate := validator.New()
	var ve validator.ValidationErrors
	err := validate.Struct(structToValidate)
	if errors.As(err, &ve) {
		return ve
	}
	return nil
}

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
