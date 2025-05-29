package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
} // @name ErrorResponse

type StatusError struct {
	err  error
	code int
}

func NewStatusError(code int, err error) *StatusError {
	return &StatusError{
		err:  err,
		code: code,
	}
}

func (e *StatusError) Error() string {
	return e.err.Error()
}

func (e *StatusError) Code() int {
	return e.code
}

func HandleError(c *gin.Context, err error) {
	statusErr, ok := err.(*StatusError)
	if ok {
		c.JSON(statusErr.code, ErrorResponse{Error: statusErr.Error()})
	} else {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}
}
