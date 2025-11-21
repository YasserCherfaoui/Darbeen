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

// SupplierBill endpoints
func (h *SupplierHandler) CreateSupplierBill(c *gin.Context) {
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

	var req supplierApp.CreateSupplierBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.supplierService.CreateSupplierBill(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Supplier bill created successfully", result)
}

func (h *SupplierHandler) GetSupplierBill(c *gin.Context) {
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

	billID, err := strconv.ParseUint(c.Param("billId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid bill id"))
		return
	}

	result, err := h.supplierService.GetSupplierBillByID(userID, uint(companyID), uint(billID))
	if err != nil {
		response.Error(c, err)
		return
	}

	// Verify supplier ID matches
	if result.SupplierID != uint(supplierID) {
		response.Error(c, errors.NewForbiddenError("bill does not belong to this supplier"))
		return
	}

	response.Success(c, http.StatusOK, result)
}

// GetSupplierBillByID gets a bill by ID only (without requiring supplier ID in path)
func (h *SupplierHandler) GetSupplierBillByID(c *gin.Context) {
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

	billID, err := strconv.ParseUint(c.Param("billId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid bill id"))
		return
	}

	result, err := h.supplierService.GetSupplierBillByID(userID, uint(companyID), uint(billID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *SupplierHandler) ListSupplierBills(c *gin.Context) {
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

	// Parse pagination parameters
	var pagination supplierApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	result, err := h.supplierService.ListSupplierBills(userID, uint(companyID), uint(supplierID), pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *SupplierHandler) UpdateSupplierBill(c *gin.Context) {
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

	billID, err := strconv.ParseUint(c.Param("billId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid bill id"))
		return
	}

	var req supplierApp.UpdateSupplierBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.supplierService.UpdateSupplierBill(userID, uint(companyID), uint(billID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Supplier bill updated successfully", result)
}

func (h *SupplierHandler) DeleteSupplierBill(c *gin.Context) {
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

	billID, err := strconv.ParseUint(c.Param("billId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid bill id"))
		return
	}

	err = h.supplierService.DeleteSupplierBill(userID, uint(companyID), uint(billID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Supplier bill deleted successfully", nil)
}

func (h *SupplierHandler) AddBillItem(c *gin.Context) {
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

	billID, err := strconv.ParseUint(c.Param("billId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid bill id"))
		return
	}

	var req supplierApp.SupplierBillItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.supplierService.AddBillItem(userID, uint(companyID), uint(billID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Bill item added successfully", result)
}

func (h *SupplierHandler) UpdateBillItem(c *gin.Context) {
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

	billID, err := strconv.ParseUint(c.Param("billId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid bill id"))
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid item id"))
		return
	}

	var req supplierApp.SupplierBillItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.supplierService.UpdateBillItem(userID, uint(companyID), uint(billID), uint(itemID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Bill item updated successfully", result)
}

func (h *SupplierHandler) RemoveBillItem(c *gin.Context) {
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

	billID, err := strconv.ParseUint(c.Param("billId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid bill id"))
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid item id"))
		return
	}

	err = h.supplierService.RemoveBillItem(userID, uint(companyID), uint(billID), uint(itemID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Bill item removed successfully", nil)
}

func (h *SupplierHandler) RecordSupplierPayment(c *gin.Context) {
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

	var req supplierApp.RecordSupplierPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.supplierService.RecordSupplierPayment(userID, uint(companyID), uint(supplierID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Supplier payment recorded successfully", result)
}

func (h *SupplierHandler) GetSupplierOutstandingBalance(c *gin.Context) {
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

	result, err := h.supplierService.GetSupplierOutstandingBalance(userID, uint(companyID), uint(supplierID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

