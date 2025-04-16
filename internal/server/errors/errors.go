package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler is a middleware that handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				})
			case gin.ErrorTypePublic:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
				})
			}
		}
	}
}

// NotFound handles 404 errors
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    http.StatusNotFound,
		Message: "Not Found",
	})
}

// MethodNotAllowed handles 405 errors
func MethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
		Code:    http.StatusMethodNotAllowed,
		Message: "Method Not Allowed",
	})
}
