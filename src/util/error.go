package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func (serviceError *ServiceError) AddDetailError(resource string, field string, error string) {
	detailError := DetailError{
		Resource: resource,
		Field:    field,
		Error:    error,
	}
	serviceError.Body.Errors = append(serviceError.Body.Errors, detailError)
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
	serviceError, ok := err.(*ServiceError)
	return serviceError, ok
}

func ErrorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) == 0 {
		return
	}
	for _, err := range c.Errors {
		Logger().Error(err)
	}
	c.JSON(http.StatusInternalServerError, "")

}

func CheckError(c *gin.Context, err error) bool {
	if err != nil {
		serviceError, ok := IsServiceError(err)
		if ok {
			c.AbortWithStatusJSON(serviceError.Build())
			return true
		}
		c.Error(err)
	}
	return false
}
