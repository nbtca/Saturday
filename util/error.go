package util

import "github.com/gin-gonic/gin"

type DetailError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ServiceError struct {
	error
	HttpStatus int
	Body       struct {
		Message string        `json:"message"`
		Errors  []DetailError `json:"errors"`
	}
}

func (error *ServiceError) setErrors(errors []DetailError) *ServiceError {
	error.Body.Errors = append(error.Body.Errors, errors...)
	return error
}

func (error *ServiceError) getCodeAndBody() (int, interface{}) {
	return error.HttpStatus, error.Body
}

func MakeServiceError(HttpStatus int) *ServiceError {
	error := &ServiceError{
		HttpStatus: HttpStatus,
	}
	return error
}

func ErrorHandler(c *gin.Context) {
	c.Next()
	// for i, err := range c.Errors {
	// }
	if len(c.Errors) > 0 {
		error := MakeServiceError(c.Writer.Status())
		error.Body.Message = c.Errors.String()
		// error.setToCtx(c)
	}
}
