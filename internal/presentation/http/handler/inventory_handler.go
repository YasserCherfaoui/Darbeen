package handler

import (
	"net/http"
	"strconv"

	inventoryApp "github.com/YasserCherfaoui/darween/internal/application/inventory"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService *inventoryApp.Service
}

func NewInventoryHandler(inventoryService *inventoryApp.Service) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req inventoryApp.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.inventoryService.CreateInventory(userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Inventory created successfully", result)
}

func (h *InventoryHandler) GetCompanyInventory(c *gin.Context) {
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

	result, err := h.inventoryService.GetInventoryByCompany(userID, uint(companyID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *InventoryHandler) GetFranchiseInventory(c *gin.Context) {
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

	result, err := h.inventoryService.GetInventoryByFranchise(userID, uint(franchiseID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *InventoryHandler) UpdateInventoryStock(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	inventoryID, err := strconv.ParseUint(c.Param("inventoryId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid inventory id"))
		return
	}

	var req inventoryApp.UpdateInventoryStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.inventoryService.UpdateInventoryStock(userID, uint(inventoryID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Inventory updated successfully", result)
}

func (h *InventoryHandler) AdjustInventoryStock(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	inventoryID, err := strconv.ParseUint(c.Param("inventoryId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid inventory id"))
		return
	}

	var req inventoryApp.AdjustInventoryStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.inventoryService.AdjustInventoryStock(userID, uint(inventoryID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Inventory adjusted successfully", result)
}

func (h *InventoryHandler) ReserveStock(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	inventoryID, err := strconv.ParseUint(c.Param("inventoryId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid inventory id"))
		return
	}

	var req inventoryApp.ReserveStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.inventoryService.ReserveStock(userID, uint(inventoryID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Stock reserved successfully", result)
}

func (h *InventoryHandler) ReleaseStock(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	inventoryID, err := strconv.ParseUint(c.Param("inventoryId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid inventory id"))
		return
	}

	var req inventoryApp.ReleaseStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.inventoryService.ReleaseStock(userID, uint(inventoryID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Stock released successfully", result)
}

func (h *InventoryHandler) GetInventoryMovements(c *gin.Context) {
	inventoryID, err := strconv.ParseUint(c.Param("inventoryId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid inventory id"))
		return
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	result, err := h.inventoryService.GetInventoryMovements(uint(inventoryID), limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}




