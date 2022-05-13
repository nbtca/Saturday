package util

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TODO refact

type DetailError struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Error    string `json:"error"`
}

type ServiceError struct {
	error
	HttpStatus int
	Body       struct {
		Message string        `json:"message"`
		Errors  []DetailError `json:"errors,omitempty"`
	}
}

func (serviceError *ServiceError) AddDetailError(resource string, field string, error string) *ServiceError {
	detailError := DetailError{
		Resource: resource,
		Field:    field,
		Error:    error,
	}
	serviceError.Body.Errors = append(serviceError.Body.Errors, detailError)
	return serviceError
}

func (error ServiceError) Build() (int, interface{}) {
	return error.HttpStatus, error.Body
}

func MakeServiceError(HttpStatus int) *ServiceError {
	error := &ServiceError{
		HttpStatus: HttpStatus,
	}
	return error
}

func (serviceError *ServiceError) SetStatus(status int) *ServiceError {
	serviceError.HttpStatus = status
	return serviceError
}

func (serviceError *ServiceError) SetMessage(message string) *ServiceError {
	serviceError.Body.Message = message
	return serviceError
}

func IsServiceError(err error) (*ServiceError, bool) {
	if err == nil {
		return nil, false
	}
	serviceError, ok := err.(*ServiceError)
	return serviceError, ok
}

// check error type
// if error is not nil, return true and handle accordingly,
// else return false
func CheckError(c *gin.Context, err error) bool {
	if err != nil {
		serviceError, ok := IsServiceError(err)
		if ok {
			c.AbortWithStatusJSON(serviceError.Build())
		} else {
			c.Error(err)
		}
		return true
	}
	return false
}

func MakeValidationError(resource string, err error) error {
	var ve validator.ValidationErrors
	serviceError := MakeServiceError(http.StatusUnprocessableEntity).SetMessage("Validation Failed")
	if err == nil {
		return serviceError
	}
	if !errors.As(err, &ve) {
		return serviceError
	}
	for _, fe := range ve {
		serviceError.AddDetailError(resource, fe.Field(), GetErrorMessage(fe))
	}
	return serviceError
}

func MakeInternalServerError() error {
	return MakeServiceError(http.StatusInternalServerError).SetMessage("Internal Server Error")
}
