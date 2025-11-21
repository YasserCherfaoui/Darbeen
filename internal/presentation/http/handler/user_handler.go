package handler

import (
	"net/http"
	"strconv"

	userApp "github.com/YasserCherfaoui/darween/internal/application/user"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *userApp.Service
}

func NewUserHandler(userService *userApp.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.userService.GetUserByID(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req userApp.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "User updated successfully", result)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	// For now, we'll require a company_id query parameter
	companyIDStr := c.Query("company_id")
	if companyIDStr == "" {
		response.Error(c, errors.NewBadRequestError("company_id query parameter is required"))
		return
	}

	companyID64, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid company_id"))
		return
	}
	companyID := uint(companyID64)

	// Note: We should verify that the requesting user has access to this company
	// This is a simplified version
	_ = userID

	result, err := h.userService.GetUsersByCompanyID(companyID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req userApp.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.userService.ChangePassword(userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, result.Message, result)
}

func (h *UserHandler) GetUserPortals(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.userService.GetUserPortals(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}
