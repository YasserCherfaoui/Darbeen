package product

import (
	"encoding/json"

	"github.com/YasserCherfaoui/darween/internal/domain/product"
)

// Product DTOs
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	SKU         string  `json:"sku" binding:"required"`
	BasePrice   float64 `json:"base_price" binding:"min=0"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	BasePrice   float64 `json:"base_price" binding:"min=0"`
	IsActive    *bool   `json:"is_active"`
}

type ProductResponse struct {
	ID          uint                     `json:"id"`
	CompanyID   uint                     `json:"company_id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	SKU         string                   `json:"sku"`
	BasePrice   float64                  `json:"base_price"`
	IsActive    bool                     `json:"is_active"`
	Variants    []ProductVariantResponse `json:"variants,omitempty"`
}

// Product Variant DTOs
type CreateProductVariantRequest struct {
	Name       string                 `json:"name" binding:"required"`
	SKU        string                 `json:"sku" binding:"required"`
	Price      float64                `json:"price" binding:"min=0"`
	Stock      int                    `json:"stock" binding:"min=0"`
	Attributes map[string]interface{} `json:"attributes"`
}

type UpdateProductVariantRequest struct {
	Name       string                 `json:"name"`
	SKU        string                 `json:"sku"`
	Price      float64                `json:"price" binding:"min=0"`
	Stock      int                    `json:"stock" binding:"min=0"`
	Attributes map[string]interface{} `json:"attributes"`
	IsActive   *bool                  `json:"is_active"`
}

type ProductVariantResponse struct {
	ID         uint                   `json:"id"`
	ProductID  uint                   `json:"product_id"`
	Name       string                 `json:"name"`
	SKU        string                 `json:"sku"`
	Price      float64                `json:"price"`
	Stock      int                    `json:"stock"`
	Attributes map[string]interface{} `json:"attributes"`
	IsActive   bool                   `json:"is_active"`
}

// Pagination DTOs
type PaginationRequest struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// Stock management DTOs
type UpdateStockRequest struct {
	Stock int `json:"stock" binding:"min=0"`
}

type AdjustStockRequest struct {
	Amount int `json:"amount"`
}

// Helper functions
func (pr *PaginationRequest) GetDefaults() {
	if pr.Page == 0 {
		pr.Page = 1
	}
	if pr.Limit == 0 {
		pr.Limit = 20
	}
}

func (pr *PaginationRequest) GetOffset() int {
	return (pr.Page - 1) * pr.Limit
}

func NewPaginatedResponse(data interface{}, total int64, page, limit int) *PaginatedResponse {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// Convert domain entities to response DTOs
func ToProductResponse(p *product.Product) *ProductResponse {
	response := &ProductResponse{
		ID:          p.ID,
		CompanyID:   p.CompanyID,
		Name:        p.Name,
		Description: p.Description,
		SKU:         p.SKU,
		BasePrice:   p.BasePrice,
		IsActive:    p.IsActive,
	}

	if len(p.Variants) > 0 {
		response.Variants = make([]ProductVariantResponse, len(p.Variants))
		for i, v := range p.Variants {
			response.Variants[i] = *ToProductVariantResponse(&v)
		}
	}

	return response
}

func ToProductVariantResponse(v *product.ProductVariant) *ProductVariantResponse {
	var attributes map[string]interface{}
	if v.Attributes != nil {
		json.Unmarshal(v.Attributes, &attributes)
	}

	return &ProductVariantResponse{
		ID:         v.ID,
		ProductID:  v.ProductID,
		Name:       v.Name,
		SKU:        v.SKU,
		Price:      v.Price,
		Stock:      v.Stock,
		Attributes: attributes,
		IsActive:   v.IsActive,
	}
}

// Convert request DTOs to domain entities
func (req *CreateProductRequest) ToProduct(companyID uint) *product.Product {
	return &product.Product{
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		BasePrice:   req.BasePrice,
		IsActive:    true,
	}
}

func (req *CreateProductVariantRequest) ToProductVariant(productID uint) *product.ProductVariant {
	var attributesJSON []byte
	if req.Attributes != nil {
		attributesJSON, _ = json.Marshal(req.Attributes)
	}

	return &product.ProductVariant{
		ProductID:  productID,
		Name:       req.Name,
		SKU:        req.SKU,
		Price:      req.Price,
		Stock:      req.Stock,
		Attributes: attributesJSON,
		IsActive:   true,
	}
}
