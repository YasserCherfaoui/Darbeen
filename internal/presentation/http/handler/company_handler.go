package handler

import (
	"net/http"
	"strconv"

	companyApp "github.com/YasserCherfaoui/darween/internal/application/company"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	companyService *companyApp.Service
}

func NewCompanyHandler(companyService *companyApp.Service) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req companyApp.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.companyService.CreateCompany(userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Company created successfully", result)
}

func (h *CompanyHandler) ListCompanies(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.companyService.GetCompaniesByUserID(userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *CompanyHandler) GetCompany(c *gin.Context) {
	companyID, err := strconv.ParseUint(c.Param("companyId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid company id"))
		return
	}

	result, err := h.companyService.GetCompanyByID(uint(companyID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
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

	var req companyApp.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.companyService.UpdateCompany(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Company updated successfully", result)
}

func (h *CompanyHandler) AddUserToCompany(c *gin.Context) {
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

	var req companyApp.AddUserToCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.companyService.AddUserToCompany(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	message := "User added to company successfully"
	if result.UserCreated {
		message = "User created and added to company successfully"
	}

	response.SuccessWithMessage(c, http.StatusOK, message, result)
}

func (h *CompanyHandler) RemoveUserFromCompany(c *gin.Context) {
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

	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid user id"))
		return
	}

	err = h.companyService.RemoveUserFromCompany(userID, uint(companyID), uint(targetUserID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "User removed from company successfully", nil)
}

func (h *CompanyHandler) ListCompanyUsers(c *gin.Context) {
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

	result, err := h.companyService.GetCompanyUsers(userID, uint(companyID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *CompanyHandler) UpdateUserRole(c *gin.Context) {
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

	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid user id"))
		return
	}

	var req companyApp.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	err = h.companyService.UpdateUserRoleInCompany(userID, uint(companyID), uint(targetUserID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "User role updated successfully", nil)
}
