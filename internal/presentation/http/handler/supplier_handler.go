package handler

import (
	"net/http"
	"strconv"

	supplierApp "github.com/YasserCherfaoui/darween/internal/application/supplier"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	supplierService *supplierApp.Service
}

func NewSupplierHandler(supplierService *supplierApp.Service) *SupplierHandler {
	return &SupplierHandler{
		supplierService: supplierService,
	}
}

// Supplier endpoints
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
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

	var req supplierApp.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.supplierService.CreateSupplier(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Supplier created successfully", result)
}

func (h *SupplierHandler) ListSuppliers(c *gin.Context) {
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

	// Parse pagination parameters
	var pagination supplierApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	result, err := h.supplierService.GetSuppliersByCompanyID(userID, uint(companyID), pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *SupplierHandler) GetSupplier(c *gin.Context) {
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

	supplierID, err := strconv.ParseUint(c.Param("supplierId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid supplier id"))
		return
	}

	result, err := h.supplierService.GetSupplierByID(userID, uint(companyID), uint(supplierID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
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

	supplierID, err := strconv.ParseUint(c.Param("supplierId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid supplier id"))
		return
	}

	var req supplierApp.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.supplierService.UpdateSupplier(userID, uint(companyID), uint(supplierID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Supplier updated successfully", result)
}

func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
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

	supplierID, err := strconv.ParseUint(c.Param("supplierId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid supplier id"))
		return
	}

	err = h.supplierService.DeleteSupplier(userID, uint(companyID), uint(supplierID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Supplier deleted successfully", nil)
}

func (h *SupplierHandler) GetSupplierProducts(c *gin.Context) {
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

	supplierID, err := strconv.ParseUint(c.Param("supplierId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid supplier id"))
		return
	}

	result, err := h.supplierService.GetSupplierProducts(userID, uint(companyID), uint(supplierID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

