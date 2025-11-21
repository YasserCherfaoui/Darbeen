package handler

import (
	"net/http"
	"strconv"

	franchiseApp "github.com/YasserCherfaoui/darween/internal/application/franchise"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type FranchiseHandler struct {
	franchiseService *franchiseApp.Service
}

func NewFranchiseHandler(franchiseService *franchiseApp.Service) *FranchiseHandler {
	return &FranchiseHandler{
		franchiseService: franchiseService,
	}
}

func (h *FranchiseHandler) CreateFranchise(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	companyID, err := strconv.ParseUint(c.Param("companyId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid company id"))
		return
	}

	var req franchiseApp.CreateFranchiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.franchiseService.CreateFranchise(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Franchise created successfully", result)
}

func (h *FranchiseHandler) ListFranchises(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	companyID, err := strconv.ParseUint(c.Param("companyId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid company id"))
		return
	}

	result, err := h.franchiseService.GetFranchisesByCompanyID(userID, uint(companyID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *FranchiseHandler) GetFranchise(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	result, err := h.franchiseService.GetFranchiseByID(userID, uint(franchiseID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *FranchiseHandler) UpdateFranchise(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	var req franchiseApp.UpdateFranchiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.franchiseService.UpdateFranchise(userID, uint(franchiseID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Franchise updated successfully", result)
}

func (h *FranchiseHandler) InitializeFranchiseInventory(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	err = h.franchiseService.InitializeFranchiseInventory(userID, uint(franchiseID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Franchise inventory initialized successfully", nil)
}

func (h *FranchiseHandler) GetFranchisePricing(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	result, err := h.franchiseService.GetFranchisePricing(userID, uint(franchiseID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *FranchiseHandler) SetFranchisePricing(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	var req franchiseApp.SetFranchisePricingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.franchiseService.SetFranchisePricing(userID, uint(franchiseID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Franchise pricing set successfully", result)
}

func (h *FranchiseHandler) DeleteFranchisePricing(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	variantID, err := strconv.ParseUint(c.Param("variantId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid variant id"))
		return
	}

	err = h.franchiseService.DeleteFranchisePricing(userID, uint(franchiseID), uint(variantID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Franchise pricing deleted successfully", nil)
}

func (h *FranchiseHandler) AddUserToFranchise(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	var req franchiseApp.AddUserToFranchiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.franchiseService.AddUserToFranchise(userID, uint(franchiseID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	message := "User added to franchise successfully"
	if result.UserCreated {
		message = "User created and added to franchise successfully"
	}

	response.SuccessWithMessage(c, http.StatusOK, message, result)
}

func (h *FranchiseHandler) RemoveUserFromFranchise(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid user id"))
		return
	}

	err = h.franchiseService.RemoveUserFromFranchise(userID, uint(franchiseID), uint(targetUserID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "User removed from franchise successfully", nil)
}

func (h *FranchiseHandler) BulkSetFranchisePricing(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	var req franchiseApp.BulkSetFranchisePricingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.franchiseService.BulkSetFranchisePricing(userID, uint(franchiseID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Franchise pricing set successfully", result)
}

func (h *FranchiseHandler) ListFranchiseUsers(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	result, err := h.franchiseService.GetFranchiseUsers(userID, uint(franchiseID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *FranchiseHandler) UpdateUserRole(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	franchiseID, err := strconv.ParseUint(c.Param("franchiseId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid franchise id"))
		return
	}

	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid user id"))
		return
	}

	var req franchiseApp.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	err = h.franchiseService.UpdateUserRoleInFranchise(userID, uint(franchiseID), uint(targetUserID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "User role updated successfully", nil)
}




