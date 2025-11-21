package handler

import (
	"net/http"
	"strconv"

	warehousebillApp "github.com/YasserCherfaoui/darween/internal/application/warehousebill"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type WarehouseBillHandler struct {
	warehouseBillService *warehousebillApp.Service
}

func NewWarehouseBillHandler(warehouseBillService *warehousebillApp.Service) *WarehouseBillHandler {
	return &WarehouseBillHandler{
		warehouseBillService: warehouseBillService,
	}
}

// CreateExitBill creates an exit bill from warehouse to franchise
func (h *WarehouseBillHandler) CreateExitBill(c *gin.Context) {
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

	var req warehousebillApp.CreateExitBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.warehouseBillService.CreateExitBill(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Exit bill created successfully", result)
}

// CreateEntryBill creates an entry bill linked to an exit bill
func (h *WarehouseBillHandler) CreateEntryBill(c *gin.Context) {
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

	var req warehousebillApp.CreateEntryBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.warehouseBillService.CreateEntryBill(userID, uint(franchiseID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Entry bill created successfully", result)
}

// VerifyEntryBill verifies received items and records discrepancies
func (h *WarehouseBillHandler) VerifyEntryBill(c *gin.Context) {
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

	billID, err := strconv.ParseUint(c.Param("billId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid bill id"))
		return
	}

	var req warehousebillApp.VerifyEntryBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.warehouseBillService.VerifyEntryBill(userID, uint(franchiseID), uint(billID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Entry bill verified successfully", result)
}

// ListWarehouseBills lists warehouse bills for a company
func (h *WarehouseBillHandler) ListWarehouseBills(c *gin.Context) {
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

	var pagination warehousebillApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	result, err := h.warehouseBillService.ListWarehouseBills(userID, uint(companyID), pagination.Page, pagination.Limit, &pagination)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// ListFranchiseWarehouseBills lists warehouse bills for a franchise
func (h *WarehouseBillHandler) ListFranchiseWarehouseBills(c *gin.Context) {
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

	var pagination warehousebillApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	result, err := h.warehouseBillService.ListFranchiseWarehouseBills(userID, uint(franchiseID), pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// GetWarehouseBill gets a single warehouse bill
func (h *WarehouseBillHandler) GetWarehouseBill(c *gin.Context) {
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

	result, err := h.warehouseBillService.GetWarehouseBill(userID, uint(companyID), uint(billID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// GetFranchiseWarehouseBill gets a single warehouse bill for a franchise
func (h *WarehouseBillHandler) GetFranchiseWarehouseBill(c *gin.Context) {
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

	billID, err := strconv.ParseUint(c.Param("billId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid bill id"))
		return
	}

	result, err := h.warehouseBillService.GetFranchiseWarehouseBill(userID, uint(franchiseID), uint(billID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// CompleteExitBill completes an exit bill
func (h *WarehouseBillHandler) CompleteExitBill(c *gin.Context) {
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

	result, err := h.warehouseBillService.CompleteExitBill(userID, uint(companyID), uint(billID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Exit bill completed successfully", result)
}

// CompleteEntryBill completes an entry bill
func (h *WarehouseBillHandler) CompleteEntryBill(c *gin.Context) {
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

	billID, err := strconv.ParseUint(c.Param("billId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid bill id"))
		return
	}

	result, err := h.warehouseBillService.CompleteEntryBill(userID, uint(franchiseID), uint(billID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Entry bill completed successfully", result)
}

// CancelWarehouseBill cancels a warehouse bill
func (h *WarehouseBillHandler) CancelWarehouseBill(c *gin.Context) {
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

	err = h.warehouseBillService.CancelWarehouseBill(userID, uint(companyID), uint(billID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Warehouse bill cancelled successfully", nil)
}

// UpdateExitBillItems updates items in an exit bill
func (h *WarehouseBillHandler) UpdateExitBillItems(c *gin.Context) {
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

	var req warehousebillApp.UpdateExitBillItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.warehouseBillService.UpdateExitBillItems(userID, uint(companyID), uint(billID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Exit bill items updated successfully", result)
}

// SearchProductsForExitBill searches for products/variants with franchise pricing
func (h *WarehouseBillHandler) SearchProductsForExitBill(c *gin.Context) {
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

	var req warehousebillApp.SearchProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.warehouseBillService.SearchProductsForExitBill(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

