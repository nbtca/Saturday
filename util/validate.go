package util

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

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

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("section", func(fl validator.FieldLevel) bool {
			section, ok := fl.Field().Interface().(string)
			if ok {
				reg := regexp.MustCompile(`^(\p{Han}{2,10})(\d{3})$`)
				if reg.MatchString(section) {
					return true
				}
			}
			return false
		})
	}
}
