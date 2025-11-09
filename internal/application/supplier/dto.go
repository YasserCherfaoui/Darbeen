package supplier

import (
	"github.com/YasserCherfaoui/darween/internal/domain/supplier"
)

// Supplier DTOs
type CreateSupplierRequest struct {
	Name          string `json:"name" binding:"required"`
	ContactPerson string `json:"contact_person"`
	Email         string `json:"email" binding:"omitempty,email"`
	Phone         string `json:"phone"`
}

type UpdateSupplierRequest struct {
	Name          string  `json:"name"`
	ContactPerson string  `json:"contact_person"`
	Email         string  `json:"email" binding:"omitempty,email"`
	Phone         string  `json:"phone"`
	IsActive      *bool   `json:"is_active"`
}

type SupplierResponse struct {
	ID            uint   `json:"id"`
	CompanyID     uint   `json:"company_id"`
	Name          string `json:"name"`
	ContactPerson string `json:"contact_person"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	IsActive      bool   `json:"is_active"`
	ProductCount  int    `json:"product_count"`
}

type SupplierWithProductsResponse struct {
	ID            uint                  `json:"id"`
	CompanyID     uint                  `json:"company_id"`
	Name          string                `json:"name"`
	ContactPerson string                `json:"contact_person"`
	Email         string                `json:"email"`
	Phone         string                `json:"phone"`
	IsActive      bool                  `json:"is_active"`
	Products      []SupplierProductInfo `json:"products"`
}

type SupplierProductInfo struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	SKU          string   `json:"sku"`
	SupplierCost *float64 `json:"supplier_cost,omitempty"`
}

// Pagination DTOs (reusing from product package pattern)
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

// Helper functions
func (pr *PaginationRequest) GetDefaults() {
	if pr.Page == 0 {
		pr.Page = 1
	}
	if pr.Limit == 0 {
		pr.Limit = 20
	}
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
func ToSupplierResponse(s *supplier.Supplier, productCount int) *SupplierResponse {
	return &SupplierResponse{
		ID:            s.ID,
		CompanyID:     s.CompanyID,
		Name:          s.Name,
		ContactPerson: s.ContactPerson,
		Email:         s.Email,
		Phone:         s.Phone,
		IsActive:      s.IsActive,
		ProductCount:  productCount,
	}
}

// Convert request DTOs to domain entities
func (req *CreateSupplierRequest) ToSupplier(companyID uint) *supplier.Supplier {
	return &supplier.Supplier{
		CompanyID:     companyID,
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Email:         req.Email,
		Phone:         req.Phone,
		IsActive:      true,
	}
}

