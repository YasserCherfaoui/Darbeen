package response

import (
	"net/http"

	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		// Unknown error, treat as internal server error
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    errors.CodeInternal,
				Message: "Internal server error",
			},
		})
		return
	}

	statusCode := getHTTPStatusCode(appErr.Code)
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}

func getHTTPStatusCode(errorCode string) int {
	switch errorCode {
	case errors.CodeNotFound:
		return http.StatusNotFound
	case errors.CodeUnauthorized:
		return http.StatusUnauthorized
	case errors.CodeForbidden:
		return http.StatusForbidden
	case errors.CodeValidation:
		return http.StatusBadRequest
	case errors.CodeConflict:
		return http.StatusConflict
	case errors.CodeBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
