package response

import (
	"encoding/json"
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
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Issues  interface{} `json:"issues,omitempty"`
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
	errorInfo := &ErrorInfo{
		Code:    appErr.Code,
		Message: appErr.Message,
	}

	// Check if the message is JSON (for validation errors with multiple issues)
	if appErr.Code == errors.CodeValidation {
		var issuesData map[string]interface{}
		if err := json.Unmarshal([]byte(appErr.Message), &issuesData); err == nil {
			// Successfully parsed as JSON, extract issues
			if issues, ok := issuesData["issues"]; ok {
				errorInfo.Issues = issues
				errorInfo.Message = "Validation failed. Please check the issues below."
			}
		}
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error:   errorInfo,
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
