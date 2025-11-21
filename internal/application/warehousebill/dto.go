package warehousebill

import (
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/warehousebill"
)

// ValidationIssue represents a single validation issue
type ValidationIssue struct {
	ItemIndex    int    `json:"item_index,omitempty"`
	VariantID    uint   `json:"variant_id,omitempty"`
	VariantSKU   string `json:"variant_sku,omitempty"`
	ProductSKU   string `json:"product_sku,omitempty"`
	ProductName  string `json:"product_name,omitempty"`
	Message      string `json:"message"`
	AvailableQty int    `json:"available_qty,omitempty"`
	RequiredQty  int    `json:"required_qty,omitempty"`
}

// ValidationErrorsResponse represents multiple validation errors
type ValidationErrorsResponse struct {
	Issues []ValidationIssue `json:"issues"`
}

// SearchProductsRequest represents a request to search products for exit bills
type SearchProductsRequest struct {
	Query       string `form:"query" binding:"required"`
	FranchiseID uint   `form:"franchise_id" binding:"required"`
	Limit       int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// ProductVariantSearchResponse represents a product variant search result with pricing
type ProductVariantSearchResponse struct {
	// Variant details
	VariantID   uint   `json:"variant_id"`
	VariantName string `json:"variant_name"`
	VariantSKU  string `json:"variant_sku"`
	
	// Product details
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductSKU  string `json:"product_sku"`
	
	// Base pricing (from product or variant)
	BaseRetailPrice    float64 `json:"base_retail_price"`
	BaseWholesalePrice float64 `json:"base_wholesale_price"`
	
	// Variant-specific pricing (if not using parent pricing)
	VariantRetailPrice    *float64 `json:"variant_retail_price,omitempty"`
	VariantWholesalePrice *float64 `json:"variant_wholesale_price,omitempty"`
	
	// Franchise pricing (if available)
	FranchiseRetailPrice    *float64 `json:"franchise_retail_price,omitempty"`
	FranchiseWholesalePrice *float64 `json:"franchise_wholesale_price,omitempty"`
	
	// Effective pricing (franchise override > variant > product base)
	EffectiveRetailPrice    float64 `json:"effective_retail_price"`
	EffectiveWholesalePrice float64 `json:"effective_wholesale_price"`
	
	// Flags
	UseParentPricing bool `json:"use_parent_pricing"`
}

// WarehouseBillItemRequest represents an item in a warehouse bill request
type WarehouseBillItemRequest struct {
	ProductVariantID uint    `json:"product_variant_id" binding:"required"`
	Quantity         int     `json:"quantity" binding:"required,min=1"`
	UnitPrice        float64 `json:"unit_price" binding:"required,min=0"`
}

// CreateExitBillRequest represents a request to create an exit bill
type CreateExitBillRequest struct {
	FranchiseID uint                      `json:"franchise_id" binding:"required"`
	Items       []WarehouseBillItemRequest `json:"items" binding:"required,min=1"`
	Notes       string                    `json:"notes"`
}

// CreateEntryBillRequest represents a request to create an entry bill
type CreateEntryBillRequest struct {
	ExitBillID uint `json:"exit_bill_id" binding:"required"`
	Notes      string `json:"notes"`
}

// VerifyEntryBillItemRequest represents an item in a verification request
type VerifyEntryBillItemRequest struct {
	ProductVariantID uint `json:"product_variant_id" binding:"required"`
	ReceivedQuantity int  `json:"received_quantity" binding:"min=0"`
}

// VerifyEntryBillRequest represents a request to verify an entry bill
type VerifyEntryBillRequest struct {
	Items []VerifyEntryBillItemRequest `json:"items" binding:"required,min=1"`
	Notes string                       `json:"notes"`
}

// UpdateExitBillItemRequest represents an item in an exit bill update request
type UpdateExitBillItemRequest struct {
	ID              *uint   `json:"id,omitempty"`              // If nil, this is a new item to add
	ProductVariantID uint   `json:"product_variant_id" binding:"required"`
	Quantity        int     `json:"quantity" binding:"required,min=1"`
	UnitPrice       float64 `json:"unit_price" binding:"required,min=0"`
}

// UpdateExitBillItemsRequest represents a request to update exit bill items
type UpdateExitBillItemsRequest struct {
	Items        []UpdateExitBillItemRequest `json:"items" binding:"required,min=1"`
	ChangeReason string                      `json:"change_reason"` // Optional reason for the change
}

// WarehouseBillItemResponse represents an item in a warehouse bill response
type WarehouseBillItemResponse struct {
	ID                uint                            `json:"id"`
	WarehouseBillID   uint                            `json:"warehouse_bill_id"`
	ProductVariantID  uint                            `json:"product_variant_id"`
	ExpectedQuantity  int                             `json:"expected_quantity"`
	ReceivedQuantity  *int                            `json:"received_quantity,omitempty"`
	Quantity          int                             `json:"quantity"`
	UnitPrice         float64                         `json:"unit_price"`
	TotalAmount       float64                         `json:"total_amount"`
	DiscrepancyType   warehousebill.DiscrepancyType   `json:"discrepancy_type"`
	DiscrepancyNotes  string                          `json:"discrepancy_notes"`
	CreatedAt         time.Time                       `json:"created_at"`
	// Product and variant details
	ProductName       *string                         `json:"product_name,omitempty"`
	VariantName       *string                         `json:"variant_name,omitempty"`
	VariantSKU       *string                         `json:"variant_sku,omitempty"`
}

// WarehouseBillResponse represents a warehouse bill response
type WarehouseBillResponse struct {
	ID                uint                            `json:"id"`
	CompanyID         uint                            `json:"company_id"`
	FranchiseID       uint                            `json:"franchise_id"`
	BillNumber        string                          `json:"bill_number"`
	BillType          warehousebill.BillType         `json:"bill_type"`
	RelatedBillID     *uint                           `json:"related_bill_id,omitempty"`
	Status            warehousebill.BillStatus        `json:"status"`
	VerificationStatus warehousebill.VerificationStatus `json:"verification_status"`
	TotalAmount       float64                         `json:"total_amount"`
	Notes             string                          `json:"notes"`
	VerifiedByID      *uint                           `json:"verified_by_id,omitempty"`
	VerifiedAt        *time.Time                      `json:"verified_at,omitempty"`
	CreatedByID       uint                            `json:"created_by_id"`
	CreatedAt         time.Time                       `json:"created_at"`
	UpdatedAt         time.Time                      `json:"updated_at"`
	Items             []WarehouseBillItemResponse     `json:"items,omitempty"`
}

// Pagination DTOs
type PaginationRequest struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
	
	// Filters
	FranchiseID *uint   `form:"franchise_id"` // Filter by franchise
	Status      *string `form:"status"`       // Filter by bill status (draft, completed, cancelled, verified)
	DateFrom    *string `form:"date_from"`    // Filter by date from (format: YYYY-MM-DD)
	DateTo      *string `form:"date_to"`      // Filter by date to (format: YYYY-MM-DD)
	BillType    *string `form:"bill_type"`    // Filter by bill type (exit, entry)
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
func ToWarehouseBillResponse(bill *warehousebill.WarehouseBill) *WarehouseBillResponse {
	response := &WarehouseBillResponse{
		ID:                bill.ID,
		CompanyID:         bill.CompanyID,
		FranchiseID:       bill.FranchiseID,
		BillNumber:        bill.BillNumber,
		BillType:          bill.BillType,
		RelatedBillID:     bill.RelatedBillID,
		Status:            bill.Status,
		VerificationStatus: bill.VerificationStatus,
		TotalAmount:       bill.TotalAmount,
		Notes:             bill.Notes,
		VerifiedByID:      bill.VerifiedByID,
		VerifiedAt:        bill.VerifiedAt,
		CreatedByID:       bill.CreatedByID,
		CreatedAt:         bill.CreatedAt,
		UpdatedAt:         bill.UpdatedAt,
	}

	if len(bill.Items) > 0 {
		response.Items = make([]WarehouseBillItemResponse, len(bill.Items))
		for i, item := range bill.Items {
			response.Items[i] = WarehouseBillItemResponse{
				ID:               item.ID,
				WarehouseBillID:   item.WarehouseBillID,
				ProductVariantID:  item.ProductVariantID,
				ExpectedQuantity:  item.ExpectedQuantity,
				ReceivedQuantity:  item.ReceivedQuantity,
				Quantity:          item.Quantity,
				UnitPrice:         item.UnitPrice,
				TotalAmount:       item.TotalAmount,
				DiscrepancyType:   item.DiscrepancyType,
				DiscrepancyNotes:  item.DiscrepancyNotes,
				CreatedAt:         item.CreatedAt,
				// ProductName, VariantName, VariantSKU will be populated by enrichWarehouseBillResponse in service layer
			}
		}
	}

	return response
}

func ToWarehouseBillItemResponse(item *warehousebill.WarehouseBillItem) *WarehouseBillItemResponse {
	return &WarehouseBillItemResponse{
		ID:               item.ID,
		WarehouseBillID:   item.WarehouseBillID,
		ProductVariantID:  item.ProductVariantID,
		ExpectedQuantity:  item.ExpectedQuantity,
		ReceivedQuantity:  item.ReceivedQuantity,
		Quantity:          item.Quantity,
		UnitPrice:         item.UnitPrice,
		TotalAmount:       item.TotalAmount,
		DiscrepancyType:   item.DiscrepancyType,
		DiscrepancyNotes:  item.DiscrepancyNotes,
		CreatedAt:         item.CreatedAt,
		// ProductName, VariantName, VariantSKU will be populated by enrichWarehouseBillResponse in service layer
	}
}

