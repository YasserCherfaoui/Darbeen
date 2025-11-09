package handler

import (
	"net/http"
	"strconv"

	productApp "github.com/YasserCherfaoui/darween/internal/application/product"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *productApp.Service
}

func NewProductHandler(productService *productApp.Service) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// Product endpoints
func (h *ProductHandler) CreateProduct(c *gin.Context) {
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

	var req productApp.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.productService.CreateProduct(userID, uint(companyID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Product created successfully", result)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
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
	var pagination productApp.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}
	pagination.GetDefaults()

	result, err := h.productService.GetProductsByCompanyID(userID, uint(companyID), pagination.Page, pagination.Limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
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

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product id"))
		return
	}

	result, err := h.productService.GetProductByID(userID, uint(companyID), uint(productID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
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

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product id"))
		return
	}

	var req productApp.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.productService.UpdateProduct(userID, uint(companyID), uint(productID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Product updated successfully", result)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
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

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product id"))
		return
	}

	err = h.productService.DeleteProduct(userID, uint(companyID), uint(productID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Product deleted successfully", nil)
}

// Product variant endpoints
func (h *ProductHandler) CreateProductVariant(c *gin.Context) {
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

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product id"))
		return
	}

	var req productApp.CreateProductVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.productService.CreateProductVariant(userID, uint(companyID), uint(productID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Product variant created successfully", result)
}

func (h *ProductHandler) ListProductVariants(c *gin.Context) {
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

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product id"))
		return
	}

	result, err := h.productService.GetProductVariantsByProductID(userID, uint(companyID), uint(productID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *ProductHandler) GetProductVariant(c *gin.Context) {
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

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product id"))
		return
	}

	variantID, err := strconv.ParseUint(c.Param("variantId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid variant id"))
		return
	}

	result, err := h.productService.GetProductVariantByID(userID, uint(companyID), uint(productID), uint(variantID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *ProductHandler) UpdateProductVariant(c *gin.Context) {
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

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product id"))
		return
	}

	variantID, err := strconv.ParseUint(c.Param("variantId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid variant id"))
		return
	}

	var req productApp.UpdateProductVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.productService.UpdateProductVariant(userID, uint(companyID), uint(productID), uint(variantID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Product variant updated successfully", result)
}

func (h *ProductHandler) DeleteProductVariant(c *gin.Context) {
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

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product id"))
		return
	}

	variantID, err := strconv.ParseUint(c.Param("variantId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid variant id"))
		return
	}

	err = h.productService.DeleteProductVariant(userID, uint(companyID), uint(productID), uint(variantID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Product variant deleted successfully", nil)
}

// BulkCreateProductVariants creates multiple product variants from attribute combinations
func (h *ProductHandler) BulkCreateProductVariants(c *gin.Context) {
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

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		response.Error(c, errors.NewBadRequestError("invalid product id"))
		return
	}

	var req productApp.BulkCreateProductVariantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.NewValidationError(err.Error()))
		return
	}

	result, err := h.productService.BulkCreateProductVariants(userID, uint(companyID), uint(productID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, http.StatusCreated, "Product variants created successfully", result)
}
