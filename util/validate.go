package util

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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

func GetBody(c *gin.Context, target interface{}) error {
	err := c.ShouldBindJSON(target)
	if err != nil {
		return err
	}
	ve := ValidationHandler(target)
	if ve == nil {
		return nil
	}
	serviceError := MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed")
	for _, fe := range ve {
		serviceError.AddDetailError(c.Request.URL.Path, fe.Field(), GetErrorMessage(fe))
	}
	return serviceError
	// c.AbortWithStatusJSON(serviceError.Build())
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
