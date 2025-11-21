package pos

import (
	"fmt"
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/pos"
)

// Customer DTOs

type CreateCustomerRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func (req *CreateCustomerRequest) ToCustomer(companyID uint) *pos.Customer {
	return &pos.Customer{
		CompanyID: companyID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		IsActive:  true,
	}
}

type UpdateCustomerRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	IsActive *bool   `json:"is_active"`
}

type CustomerResponse struct {
	ID             uint      `json:"id"`
	CompanyID      uint      `json:"company_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	TotalPurchases float64   `json:"total_purchases"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ToCustomerResponse(customer *pos.Customer) *CustomerResponse {
	return &CustomerResponse{
		ID:             customer.ID,
		CompanyID:      customer.CompanyID,
		Name:           customer.Name,
		Email:          customer.Email,
		Phone:          customer.Phone,
		Address:        customer.Address,
		TotalPurchases: customer.TotalPurchases,
		IsActive:       customer.IsActive,
		CreatedAt:      customer.CreatedAt,
		UpdatedAt:      customer.UpdatedAt,
	}
}

// Sale DTOs

type SaleItemRequest struct {
	ProductVariantID uint    `json:"product_variant_id" binding:"required"`
	Quantity         int     `json:"quantity" binding:"required,min=1"`
	UnitPrice        float64 `json:"unit_price" binding:"required,min=0"`
	DiscountAmount   float64 `json:"discount_amount"`
}

type CreateSaleRequest struct {
	FranchiseID    *uint             `json:"franchise_id"`
	CustomerID     *uint             `json:"customer_id"`
	Items          []SaleItemRequest `json:"items" binding:"required,min=1"`
	TaxAmount      float64           `json:"tax_amount"`
	DiscountAmount float64           `json:"discount_amount"`
	Notes          string            `json:"notes"`
}

type SaleItemResponse struct {
	ID               uint      `json:"id"`
	SaleID           uint      `json:"sale_id"`
	ProductVariantID uint      `json:"product_variant_id"`
	Quantity         int       `json:"quantity"`
	UnitPrice        float64   `json:"unit_price"`
	DiscountAmount   float64   `json:"discount_amount"`
	SubTotal         float64   `json:"sub_total"`
	TotalAmount      float64   `json:"total_amount"`
	CreatedAt        time.Time `json:"created_at"`
	// Product and variant details
	ProductName      *string   `json:"product_name,omitempty"`
	VariantName      *string   `json:"variant_name,omitempty"`
	VariantSKU       *string   `json:"variant_sku,omitempty"`
}

func ToSaleItemResponse(item *pos.SaleItem, productName, variantName, variantSKU *string) *SaleItemResponse {
	return &SaleItemResponse{
		ID:               item.ID,
		SaleID:           item.SaleID,
		ProductVariantID: item.ProductVariantID,
		Quantity:         item.Quantity,
		UnitPrice:        item.UnitPrice,
		DiscountAmount:   item.DiscountAmount,
		SubTotal:         item.SubTotal,
		TotalAmount:      item.TotalAmount,
		CreatedAt:        item.CreatedAt,
		ProductName:      productName,
		VariantName:      variantName,
		VariantSKU:       variantSKU,
	}
}

type SaleResponse struct {
	ID             uint               `json:"id"`
	CompanyID      uint               `json:"company_id"`
	FranchiseID    *uint              `json:"franchise_id"`
	CustomerID     *uint              `json:"customer_id"`
	ReceiptNumber  string             `json:"receipt_number"`
	SubTotal       float64            `json:"sub_total"`
	TaxAmount      float64            `json:"tax_amount"`
	DiscountAmount float64            `json:"discount_amount"`
	TotalAmount    float64            `json:"total_amount"`
	PaymentStatus  pos.PaymentStatus  `json:"payment_status"`
	SaleStatus     pos.SaleStatus     `json:"sale_status"`
	Notes          string             `json:"notes"`
	CreatedByID    uint               `json:"created_by_id"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	Items          []SaleItemResponse `json:"items,omitempty"`
	Payments       []PaymentResponse  `json:"payments,omitempty"`
	Customer       *CustomerResponse  `json:"customer,omitempty"`
}

func ToSaleResponse(sale *pos.Sale) *SaleResponse {
	return ToSaleResponseWithDetails(sale, nil)
}

func ToSaleResponseWithDetails(sale *pos.Sale, itemDetails map[uint]struct {
	productName string
	variantName string
	variantSKU  string
}) *SaleResponse {
	response := &SaleResponse{
		ID:             sale.ID,
		CompanyID:      sale.CompanyID,
		FranchiseID:    sale.FranchiseID,
		CustomerID:     sale.CustomerID,
		ReceiptNumber:  sale.ReceiptNumber,
		SubTotal:       sale.SubTotal,
		TaxAmount:      sale.TaxAmount,
		DiscountAmount: sale.DiscountAmount,
		TotalAmount:    sale.TotalAmount,
		PaymentStatus:  sale.PaymentStatus,
		SaleStatus:     sale.SaleStatus,
		Notes:          sale.Notes,
		CreatedByID:    sale.CreatedByID,
		CreatedAt:      sale.CreatedAt,
		UpdatedAt:      sale.UpdatedAt,
	}

	if len(sale.Items) > 0 {
		response.Items = make([]SaleItemResponse, len(sale.Items))
		for i, item := range sale.Items {
			var productName, variantName, variantSKU *string
			if itemDetails != nil {
				if details, exists := itemDetails[item.ProductVariantID]; exists {
					if details.productName != "" {
						productName = &details.productName
					}
					if details.variantName != "" {
						variantName = &details.variantName
					}
					if details.variantSKU != "" {
						variantSKU = &details.variantSKU
					}
				}
			}
			response.Items[i] = *ToSaleItemResponse(&item, productName, variantName, variantSKU)
		}
	}

	if len(sale.Payments) > 0 {
		response.Payments = make([]PaymentResponse, len(sale.Payments))
		for i, payment := range sale.Payments {
			response.Payments[i] = *ToPaymentResponse(&payment)
		}
	}

	if sale.Customer != nil {
		response.Customer = ToCustomerResponse(sale.Customer)
	}

	return response
}

// Payment DTOs

type AddPaymentRequest struct {
	PaymentMethod pos.PaymentMethod `json:"payment_method" binding:"required"`
	Amount        float64           `json:"amount" binding:"required,min=0"`
	Reference     string            `json:"reference"`
	Notes         string            `json:"notes"`
}

type PaymentResponse struct {
	ID            uint                         `json:"id"`
	SaleID        uint                         `json:"sale_id"`
	PaymentMethod pos.PaymentMethod            `json:"payment_method"`
	Amount        float64                      `json:"amount"`
	PaymentStatus pos.PaymentTransactionStatus `json:"payment_status"`
	Reference     string                       `json:"reference"`
	Notes         string                       `json:"notes"`
	CreatedAt     time.Time                    `json:"created_at"`
}

func ToPaymentResponse(payment *pos.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:            payment.ID,
		SaleID:        payment.SaleID,
		PaymentMethod: payment.PaymentMethod,
		Amount:        payment.Amount,
		PaymentStatus: payment.PaymentStatus,
		Reference:     payment.Reference,
		Notes:         payment.Notes,
		CreatedAt:     payment.CreatedAt,
	}
}

// Cash Drawer DTOs

type OpenCashDrawerRequest struct {
	FranchiseID    *uint   `json:"franchise_id"`
	OpeningBalance float64 `json:"opening_balance" binding:"required,min=0"`
	Notes          string  `json:"notes"`
}

type CloseCashDrawerRequest struct {
	ClosingBalance float64 `json:"closing_balance" binding:"required,min=0"`
	Notes          string  `json:"notes"`
}

type CashDrawerTransactionResponse struct {
	ID              uint                          `json:"id"`
	CashDrawerID    uint                          `json:"cash_drawer_id"`
	TransactionType pos.CashDrawerTransactionType `json:"transaction_type"`
	Amount          float64                       `json:"amount"`
	SaleID          *uint                         `json:"sale_id"`
	Notes           string                        `json:"notes"`
	CreatedAt       time.Time                     `json:"created_at"`
}

func ToCashDrawerTransactionResponse(transaction *pos.CashDrawerTransaction) *CashDrawerTransactionResponse {
	return &CashDrawerTransactionResponse{
		ID:              transaction.ID,
		CashDrawerID:    transaction.CashDrawerID,
		TransactionType: transaction.TransactionType,
		Amount:          transaction.Amount,
		SaleID:          transaction.SaleID,
		Notes:           transaction.Notes,
		CreatedAt:       transaction.CreatedAt,
	}
}

type CashDrawerResponse struct {
	ID              uint                            `json:"id"`
	CompanyID       *uint                           `json:"company_id"`
	FranchiseID     *uint                           `json:"franchise_id"`
	OpeningBalance  float64                         `json:"opening_balance"`
	ClosingBalance  *float64                        `json:"closing_balance"`
	ExpectedBalance *float64                        `json:"expected_balance"`
	Difference      *float64                        `json:"difference"`
	Status          pos.CashDrawerStatus            `json:"status"`
	OpenedByID      uint                            `json:"opened_by_id"`
	ClosedByID      *uint                           `json:"closed_by_id"`
	OpenedAt        time.Time                       `json:"opened_at"`
	ClosedAt        *time.Time                      `json:"closed_at"`
	Notes           string                          `json:"notes"`
	Transactions    []CashDrawerTransactionResponse `json:"transactions,omitempty"`
	CreatedAt       time.Time                       `json:"created_at"`
	UpdatedAt       time.Time                       `json:"updated_at"`
}

func ToCashDrawerResponse(drawer *pos.CashDrawer) *CashDrawerResponse {
	response := &CashDrawerResponse{
		ID:              drawer.ID,
		CompanyID:       drawer.CompanyID,
		FranchiseID:     drawer.FranchiseID,
		OpeningBalance:  drawer.OpeningBalance,
		ClosingBalance:  drawer.ClosingBalance,
		ExpectedBalance: drawer.ExpectedBalance,
		Difference:      drawer.Difference,
		Status:          drawer.Status,
		OpenedByID:      drawer.OpenedByID,
		ClosedByID:      drawer.ClosedByID,
		OpenedAt:        drawer.OpenedAt,
		ClosedAt:        drawer.ClosedAt,
		Notes:           drawer.Notes,
		CreatedAt:       drawer.CreatedAt,
		UpdatedAt:       drawer.UpdatedAt,
	}

	if len(drawer.Transactions) > 0 {
		response.Transactions = make([]CashDrawerTransactionResponse, len(drawer.Transactions))
		for i, transaction := range drawer.Transactions {
			response.Transactions[i] = *ToCashDrawerTransactionResponse(&transaction)
		}
	}

	return response
}

// Refund DTOs

type ProcessRefundRequest struct {
	RefundAmount float64           `json:"refund_amount" binding:"required,min=0"`
	Reason       string            `json:"reason" binding:"required"`
	RefundMethod pos.PaymentMethod `json:"refund_method" binding:"required"`
}

type RefundResponse struct {
	ID             uint              `json:"id"`
	OriginalSaleID uint              `json:"original_sale_id"`
	RefundAmount   float64           `json:"refund_amount"`
	Reason         string            `json:"reason"`
	RefundMethod   pos.PaymentMethod `json:"refund_method"`
	RefundStatus   pos.RefundStatus  `json:"refund_status"`
	ProcessedByID  uint              `json:"processed_by_id"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	OriginalSale   *SaleResponse     `json:"original_sale,omitempty"`
}

func ToRefundResponse(refund *pos.Refund) *RefundResponse {
	response := &RefundResponse{
		ID:             refund.ID,
		OriginalSaleID: refund.OriginalSaleID,
		RefundAmount:   refund.RefundAmount,
		Reason:         refund.Reason,
		RefundMethod:   refund.RefundMethod,
		RefundStatus:   refund.RefundStatus,
		ProcessedByID:  refund.ProcessedByID,
		CreatedAt:      refund.CreatedAt,
		UpdatedAt:      refund.UpdatedAt,
	}

	if refund.OriginalSale != nil {
		response.OriginalSale = ToSaleResponse(refund.OriginalSale)
	}

	return response
}

// Sales Report DTOs

type SalesReportRequest struct {
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	FranchiseID *uint     `json:"franchise_id"`
}

type SalesReportResponse struct {
	TotalSales        int64              `json:"total_sales"`
	TotalRevenue      float64            `json:"total_revenue"`
	TotalCash         float64            `json:"total_cash"`
	TotalCard         float64            `json:"total_card"`
	TotalRefunded     float64            `json:"total_refunded"`
	AverageOrderValue float64            `json:"average_order_value"`
	SalesByDate       map[string]float64 `json:"sales_by_date"`
}

func ToSalesReportResponse(data *pos.SalesReportData) *SalesReportResponse {
	return &SalesReportResponse{
		TotalSales:        data.TotalSales,
		TotalRevenue:      data.TotalRevenue,
		TotalCash:         data.TotalCash,
		TotalCard:         data.TotalCard,
		TotalRefunded:     data.TotalRefunded,
		AverageOrderValue: data.AverageOrderValue,
		SalesByDate:       data.SalesByDate,
	}
}

// Pagination

type PaginationRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (p *PaginationRequest) GetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

func NewPaginatedResponse(data interface{}, total int64, page, limit int) *PaginatedResponse {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// Receipt DTOs

type ReceiptResponse struct {
	ReceiptNumber  string             `json:"receipt_number"`
	Date           time.Time          `json:"date"`
	CompanyName    string             `json:"company_name"`
	FranchiseName  string             `json:"franchise_name,omitempty"`
	CustomerName   string             `json:"customer_name,omitempty"`
	Items          []SaleItemResponse `json:"items"`
	SubTotal       float64            `json:"sub_total"`
	TaxAmount      float64            `json:"tax_amount"`
	DiscountAmount float64            `json:"discount_amount"`
	TotalAmount    float64            `json:"total_amount"`
	Payments       []PaymentResponse  `json:"payments"`
}

// Product Search DTOs

type SearchProductsRequest struct {
	Query       string `form:"query" binding:"required"`
	FranchiseID *uint  `form:"franchise_id"`
	Limit       int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// ProductVariantSearchResponse represents a product variant search result with retail pricing for POS
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
	
	// Effective pricing (franchise override > variant > product base) - for sales, we use retail price
	EffectiveRetailPrice    float64 `json:"effective_retail_price"`
	EffectiveWholesalePrice float64 `json:"effective_wholesale_price"`
	
	// Flags
	UseParentPricing bool `json:"use_parent_pricing"`
}

// Generate receipt number
func GenerateReceiptNumber(companyID uint, saleID uint) string {
	timestamp := time.Now().Format("20060102")
	return fmt.Sprintf("RCP-%d-%s-%d", companyID, timestamp, saleID)
}
