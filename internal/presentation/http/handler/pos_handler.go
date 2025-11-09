package handler

import (
	"net/http"
	"strconv"

	posApp "github.com/YasserCherfaoui/darween/internal/application/pos"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type POSHandler struct {
	posService *posApp.Service
}

func NewPOSHandler(posService *posApp.Service) *POSHandler {
	return &POSHandler{
		posService: posService,
	}
}

// Customer endpoints

func (h *POSHandler) CreateCustomer(c *gin.Context) {
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

	var req posApp.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.posService.CreateCustomer(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Customer created successfully", result)
}

func (h *POSHandler) UpdateCustomer(c *gin.Context) {
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

	customerID, err := strconv.ParseUint(c.Param("customerId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid customer id"))
		return
	}

	var req posApp.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.posService.UpdateCustomer(userID, uint(companyID), uint(customerID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Customer updated successfully", result)
}

func (h *POSHandler) GetCustomer(c *gin.Context) {
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

	customerID, err := strconv.ParseUint(c.Param("customerId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid customer id"))
		return
	}

	result, err := h.posService.GetCustomerByID(userID, uint(companyID), uint(customerID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *POSHandler) ListCustomers(c *gin.Context) {
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

	var pagination posApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	result, err := h.posService.ListCustomers(userID, uint(companyID), pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *POSHandler) DeleteCustomer(c *gin.Context) {
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

	customerID, err := strconv.ParseUint(c.Param("customerId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid customer id"))
		return
	}

	err = h.posService.DeleteCustomer(userID, uint(companyID), uint(customerID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Customer deleted successfully", nil)
}

// Sale endpoints

func (h *POSHandler) CreateSale(c *gin.Context) {
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

	var req posApp.CreateSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.posService.CreateSale(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Sale created successfully", result)
}

func (h *POSHandler) GetSale(c *gin.Context) {
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

	saleID, err := strconv.ParseUint(c.Param("saleId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid sale id"))
		return
	}

	result, err := h.posService.GetSaleByID(userID, uint(companyID), uint(saleID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *POSHandler) ListSales(c *gin.Context) {
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

	var pagination posApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	// Check for franchise filter
	var franchiseID *uint
	if franchiseIDStr := c.Query("franchise_id"); franchiseIDStr != "" {
		fID, err := strconv.ParseUint(franchiseIDStr, 10, 32)
		if err == nil {
			fIDUint := uint(fID)
			franchiseID = &fIDUint
		}
	}

	result, err := h.posService.ListSales(userID, uint(companyID), franchiseID, pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *POSHandler) ListFranchiseSales(c *gin.Context) {
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

	var pagination posApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	// Get franchise to find company ID
	// This is a simplified approach - in production, you'd fetch the franchise first
	fIDUint := uint(franchiseID)
	result, err := h.posService.ListSales(userID, 0, &fIDUint, pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// Payment endpoints

func (h *POSHandler) AddPayment(c *gin.Context) {
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

	saleID, err := strconv.ParseUint(c.Param("saleId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid sale id"))
		return
	}

	var req posApp.AddPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.posService.AddPaymentToSale(userID, uint(companyID), uint(saleID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Payment added successfully", result)
}

// Refund endpoints

func (h *POSHandler) ProcessRefund(c *gin.Context) {
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

	saleID, err := strconv.ParseUint(c.Param("saleId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid sale id"))
		return
	}

	var req posApp.ProcessRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.posService.ProcessRefund(userID, uint(companyID), uint(saleID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Refund processed successfully", result)
}

func (h *POSHandler) ListRefunds(c *gin.Context) {
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

	var pagination posApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	// Check for franchise filter
	var franchiseID *uint
	if franchiseIDStr := c.Query("franchise_id"); franchiseIDStr != "" {
		fID, err := strconv.ParseUint(franchiseIDStr, 10, 32)
		if err == nil {
			fIDUint := uint(fID)
			franchiseID = &fIDUint
		}
	}

	result, err := h.posService.ListRefunds(userID, uint(companyID), franchiseID, pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *POSHandler) ListFranchiseRefunds(c *gin.Context) {
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

	var pagination posApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	fIDUint := uint(franchiseID)
	result, err := h.posService.ListRefunds(userID, 0, &fIDUint, pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// Cash Drawer endpoints

func (h *POSHandler) OpenCashDrawer(c *gin.Context) {
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

	var req posApp.OpenCashDrawerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.posService.OpenCashDrawer(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Cash drawer opened successfully", result)
}

func (h *POSHandler) CloseCashDrawer(c *gin.Context) {
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

	drawerID, err := strconv.ParseUint(c.Param("drawerId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid drawer id"))
		return
	}

	var req posApp.CloseCashDrawerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.posService.CloseCashDrawer(userID, uint(companyID), uint(drawerID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Cash drawer closed successfully", result)
}

func (h *POSHandler) GetActiveCashDrawer(c *gin.Context) {
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

	// Check for franchise filter
	var franchiseID *uint
	if franchiseIDStr := c.Query("franchise_id"); franchiseIDStr != "" {
		fID, err := strconv.ParseUint(franchiseIDStr, 10, 32)
		if err == nil {
			fIDUint := uint(fID)
			franchiseID = &fIDUint
		}
	}

	result, err := h.posService.GetActiveCashDrawer(userID, uint(companyID), franchiseID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *POSHandler) GetActiveFranchiseCashDrawer(c *gin.Context) {
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

	fIDUint := uint(franchiseID)
	result, err := h.posService.GetActiveCashDrawer(userID, 0, &fIDUint)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *POSHandler) ListCashDrawers(c *gin.Context) {
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

	var pagination posApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	// Check for franchise filter
	var franchiseID *uint
	if franchiseIDStr := c.Query("franchise_id"); franchiseIDStr != "" {
		fID, err := strconv.ParseUint(franchiseIDStr, 10, 32)
		if err == nil {
			fIDUint := uint(fID)
			franchiseID = &fIDUint
		}
	}

	result, err := h.posService.ListCashDrawers(userID, uint(companyID), franchiseID, pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *POSHandler) ListFranchiseCashDrawers(c *gin.Context) {
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

	var pagination posApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	fIDUint := uint(franchiseID)
	result, err := h.posService.ListCashDrawers(userID, 0, &fIDUint, pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// Report endpoints

func (h *POSHandler) GetSalesReport(c *gin.Context) {
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

	var req posApp.SalesReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.posService.GetSalesReport(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *POSHandler) GetFranchiseSalesReport(c *gin.Context) {
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

	var req posApp.SalesReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	fIDUint := uint(franchiseID)
	req.FranchiseID = &fIDUint

	result, err := h.posService.GetSalesReport(userID, 0, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

