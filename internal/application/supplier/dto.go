package supplier

import (
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/supplier"
)

// Supplier DTOs
type CreateSupplierRequest struct {
	Name          string `json:"name" binding:"required"`
	ContactPerson string `json:"contact_person"`
	Email         string `json:"email" binding:"omitempty,email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
}

type UpdateSupplierRequest struct {
	Name          string  `json:"name"`
	ContactPerson string  `json:"contact_person"`
	Email         string  `json:"email" binding:"omitempty,email"`
	Phone         string  `json:"phone"`
	Address       string  `json:"address"`
	IsActive      *bool   `json:"is_active"`
}

type SupplierResponse struct {
	ID            uint   `json:"id"`
	CompanyID     uint   `json:"company_id"`
	Name          string `json:"name"`
	ContactPerson string `json:"contact_person"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
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
	Address       string                `json:"address"`
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
		Address:       s.Address,
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
		Address:       req.Address,
		IsActive:      true,
	}
}

// SupplierBill DTOs
type SupplierBillItemRequest struct {
	ProductVariantID uint    `json:"product_variant_id" binding:"required"`
	Quantity         int     `json:"quantity" binding:"required,min=1"`
	UnitCost         float64 `json:"unit_cost" binding:"required,min=0"`
}

type CreateSupplierBillRequest struct {
	SupplierID uint                     `json:"supplier_id" binding:"required"`
	Items      []SupplierBillItemRequest `json:"items" binding:"required,min=1"`
	PaidAmount float64                  `json:"paid_amount" binding:"min=0"` // Initial paid amount (can be 0)
	Notes      string                   `json:"notes"`
}

type UpdateSupplierBillRequest struct {
	Items      []SupplierBillItemRequest `json:"items"`
	Notes      string                    `json:"notes"`
	BillStatus *string                   `json:"bill_status"` // "draft", "completed", "cancelled"
}

type SupplierBillItemResponse struct {
	ID               uint      `json:"id"`
	SupplierBillID   uint      `json:"supplier_bill_id"`
	ProductVariantID uint      `json:"product_variant_id"`
	Quantity         int       `json:"quantity"`
	UnitCost         float64   `json:"unit_cost"`
	TotalCost        float64   `json:"total_cost"`
	CreatedAt        time.Time `json:"created_at"`
	// Product and variant details
	ProductName      *string   `json:"product_name,omitempty"`
	VariantName      *string   `json:"variant_name,omitempty"`
	VariantSKU       *string   `json:"variant_sku,omitempty"`
}

type SupplierBillResponse struct {
	ID             uint                      `json:"id"`
	CompanyID      uint                      `json:"company_id"`
	SupplierID     uint                      `json:"supplier_id"`
	BillNumber     string                    `json:"bill_number"`
	TotalAmount    float64                   `json:"total_amount"`
	PaidAmount     float64                   `json:"paid_amount"`
	PendingAmount  float64                   `json:"pending_amount"`
	PaymentStatus  supplier.PaymentStatus    `json:"payment_status"`
	BillStatus     supplier.BillStatus       `json:"bill_status"`
	Notes          string                    `json:"notes"`
	CreatedByID    uint                      `json:"created_by_id"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
	Items          []SupplierBillItemResponse `json:"items,omitempty"`
	Supplier       *SupplierResponse          `json:"supplier,omitempty"`
	Payments       []SupplierPaymentResponse  `json:"payments,omitempty"`
}

func ToSupplierBillResponse(bill *supplier.SupplierBill) *SupplierBillResponse {
	response := &SupplierBillResponse{
		ID:            bill.ID,
		CompanyID:     bill.CompanyID,
		SupplierID:    bill.SupplierID,
		BillNumber:    bill.BillNumber,
		TotalAmount:   bill.TotalAmount,
		PaidAmount:    bill.PaidAmount,
		PendingAmount: bill.PendingAmount,
		PaymentStatus: bill.PaymentStatus,
		BillStatus:    bill.BillStatus,
		Notes:         bill.Notes,
		CreatedByID:   bill.CreatedByID,
		CreatedAt:     bill.CreatedAt,
		UpdatedAt:     bill.UpdatedAt,
	}

	if len(bill.Items) > 0 {
		response.Items = make([]SupplierBillItemResponse, len(bill.Items))
		for i, item := range bill.Items {
			var productName, variantName, variantSKU *string
			if item.ProductVariant != nil {
				if item.ProductVariant.Product != nil {
					productName = &item.ProductVariant.Product.Name
				}
				variantName = &item.ProductVariant.Name
				variantSKU = &item.ProductVariant.SKU
			}
			response.Items[i] = SupplierBillItemResponse{
				ID:               item.ID,
				SupplierBillID:   item.SupplierBillID,
				ProductVariantID: item.ProductVariantID,
				Quantity:         item.Quantity,
				UnitCost:         item.UnitCost,
				TotalCost:        item.TotalCost,
				CreatedAt:        item.CreatedAt,
				ProductName:      productName,
				VariantName:      variantName,
				VariantSKU:       variantSKU,
			}
		}
	}

	if bill.Supplier != nil {
		response.Supplier = ToSupplierResponse(bill.Supplier, 0)
	}

	return response
}

func ToSupplierBillItemResponse(item *supplier.SupplierBillItem) *SupplierBillItemResponse {
	var productName, variantName, variantSKU *string
	if item.ProductVariant != nil {
		if item.ProductVariant.Product != nil {
			productName = &item.ProductVariant.Product.Name
		}
		variantName = &item.ProductVariant.Name
		variantSKU = &item.ProductVariant.SKU
	}
	return &SupplierBillItemResponse{
		ID:               item.ID,
		SupplierBillID:   item.SupplierBillID,
		ProductVariantID: item.ProductVariantID,
		Quantity:         item.Quantity,
		UnitCost:         item.UnitCost,
		TotalCost:        item.TotalCost,
		CreatedAt:        item.CreatedAt,
		ProductName:      productName,
		VariantName:      variantName,
		VariantSKU:       variantSKU,
	}
}

// SupplierPayment DTOs
type RecordSupplierPaymentRequest struct {
	SupplierID    uint                  `json:"supplier_id" binding:"required"`
	Amount        float64               `json:"amount" binding:"required,min=0"`
	PaymentMethod supplier.PaymentMethod `json:"payment_method" binding:"required"`
	Reference     string                `json:"reference"`
	Notes         string                `json:"notes"`
}

type SupplierPaymentDistributionResponse struct {
	ID                uint      `json:"id"`
	SupplierPaymentID uint      `json:"supplier_payment_id"`
	SupplierBillID    uint      `json:"supplier_bill_id"`
	Amount            float64   `json:"amount"`
	CreatedAt         time.Time `json:"created_at"`
}

type SupplierPaymentResponse struct {
	ID            uint                                 `json:"id"`
	SupplierID    uint                                 `json:"supplier_id"`
	CompanyID     uint                                 `json:"company_id"`
	Amount        float64                              `json:"amount"`
	PaymentMethod supplier.PaymentMethod               `json:"payment_method"`
	PaymentStatus supplier.PaymentTransactionStatus    `json:"payment_status"`
	Reference     string                               `json:"reference"`
	Notes         string                               `json:"notes"`
	CreatedByID   uint                                 `json:"created_by_id"`
	CreatedAt     time.Time                            `json:"created_at"`
	Distributions []SupplierPaymentDistributionResponse `json:"distributions,omitempty"`
}

func ToSupplierPaymentResponse(payment *supplier.SupplierPayment) *SupplierPaymentResponse {
	response := &SupplierPaymentResponse{
		ID:            payment.ID,
		SupplierID:    payment.SupplierID,
		CompanyID:     payment.CompanyID,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		PaymentStatus: payment.PaymentStatus,
		Reference:     payment.Reference,
		Notes:         payment.Notes,
		CreatedByID:   payment.CreatedByID,
		CreatedAt:     payment.CreatedAt,
	}

	if len(payment.Distributions) > 0 {
		response.Distributions = make([]SupplierPaymentDistributionResponse, len(payment.Distributions))
		for i, dist := range payment.Distributions {
			response.Distributions[i] = SupplierPaymentDistributionResponse{
				ID:                dist.ID,
				SupplierPaymentID: dist.SupplierPaymentID,
				SupplierBillID:    dist.SupplierBillID,
				Amount:            dist.Amount,
				CreatedAt:         dist.CreatedAt,
			}
		}
	}

	return response
}

// Outstanding Balance DTO
type SupplierOutstandingBalanceResponse struct {
	SupplierID      uint    `json:"supplier_id"`
	CompanyID       uint    `json:"company_id"`
	OutstandingAmount float64 `json:"outstanding_amount"`
}

