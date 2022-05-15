package util

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

var section validator.Func = func(fl validator.FieldLevel) bool {
	section, ok := fl.Field().Interface().(string)
	if ok {
		reg := regexp.MustCompile(`^(\p{Han}{2,10})(\d{3})$`)
		if reg.MatchString(section) {
			return true
		}
	}
	return false
}

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("section", section)
	}
}
