package pos

import "time"

// Customer represents a customer entity
type Customer struct {
	ID             uint   `gorm:"primaryKey"`
	CompanyID      uint   `gorm:"not null;index"`
	Name           string `gorm:"not null"`
	Email          string `gorm:"index"`
	Phone          string
	Address        string
	TotalPurchases float64 `gorm:"type:decimal(10,2);default:0"`
	IsActive       bool    `gorm:"default:true"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (Customer) TableName() string {
	return "customers"
}

// IsValid validates the customer entity
func (c *Customer) IsValid() bool {
	return c.Name != "" && c.CompanyID > 0
}

// AddPurchase adds to the total purchases amount
func (c *Customer) AddPurchase(amount float64) {
	if amount > 0 {
		c.TotalPurchases += amount
	}
}

// SaleStatus represents the status of a sale
type SaleStatus string

const (
	SaleStatusDraft     SaleStatus = "draft"
	SaleStatusCompleted SaleStatus = "completed"
	SaleStatusCancelled SaleStatus = "cancelled"
	SaleStatusRefunded  SaleStatus = "refunded"
)

func (s SaleStatus) IsValid() bool {
	switch s {
	case SaleStatusDraft, SaleStatusCompleted, SaleStatusCancelled, SaleStatusRefunded:
		return true
	}
	return false
}

// PaymentStatus represents the payment status of a sale
type PaymentStatus string

const (
	PaymentStatusUnpaid       PaymentStatus = "unpaid"
	PaymentStatusPartiallyPaid PaymentStatus = "partially_paid"
	PaymentStatusPaid         PaymentStatus = "paid"
	PaymentStatusRefunded     PaymentStatus = "refunded"
)

func (ps PaymentStatus) IsValid() bool {
	switch ps {
	case PaymentStatusUnpaid, PaymentStatusPartiallyPaid, PaymentStatusPaid, PaymentStatusRefunded:
		return true
	}
	return false
}

// Sale represents a sales transaction
type Sale struct {
	ID             uint          `gorm:"primaryKey"`
	CompanyID      uint          `gorm:"not null;index"`
	FranchiseID    *uint         `gorm:"index"`
	CustomerID     *uint         `gorm:"index"`
	ReceiptNumber  string        `gorm:"uniqueIndex;not null"`
	SubTotal       float64       `gorm:"type:decimal(10,2);not null"`
	TaxAmount      float64       `gorm:"type:decimal(10,2);default:0"`
	DiscountAmount float64       `gorm:"type:decimal(10,2);default:0"`
	TotalAmount    float64       `gorm:"type:decimal(10,2);not null"`
	PaymentStatus  PaymentStatus `gorm:"type:varchar(50);not null;default:'unpaid'"`
	SaleStatus     SaleStatus    `gorm:"type:varchar(50);not null;default:'draft'"`
	Notes          string        `gorm:"type:text"`
	CreatedByID    uint          `gorm:"not null;index"`
	CreatedAt      time.Time     `gorm:"index"`
	UpdatedAt      time.Time

	// Relationships
	Items    []SaleItem `gorm:"foreignKey:SaleID"`
	Payments []Payment  `gorm:"foreignKey:SaleID"`
	Customer *Customer  `gorm:"foreignKey:CustomerID"`
}

func (Sale) TableName() string {
	return "sales"
}

// IsValid validates the sale entity
func (s *Sale) IsValid() bool {
	return s.CompanyID > 0 && s.CreatedByID > 0 && s.TotalAmount >= 0 &&
		s.PaymentStatus.IsValid() && s.SaleStatus.IsValid()
}

// CalculateTotals recalculates all totals based on items
func (s *Sale) CalculateTotals() {
	s.SubTotal = 0
	for _, item := range s.Items {
		s.SubTotal += item.TotalAmount
	}
	s.TotalAmount = s.SubTotal + s.TaxAmount - s.DiscountAmount
}

// UpdatePaymentStatus updates the payment status based on total payments
func (s *Sale) UpdatePaymentStatus(totalPaid float64) {
	if totalPaid == 0 {
		s.PaymentStatus = PaymentStatusUnpaid
	} else if totalPaid >= s.TotalAmount {
		s.PaymentStatus = PaymentStatusPaid
	} else {
		s.PaymentStatus = PaymentStatusPartiallyPaid
	}
}

// CanBeRefunded checks if the sale can be refunded
func (s *Sale) CanBeRefunded() bool {
	return s.SaleStatus == SaleStatusCompleted && s.PaymentStatus == PaymentStatusPaid
}

// Complete marks the sale as completed
func (s *Sale) Complete() {
	s.SaleStatus = SaleStatusCompleted
}

// Cancel marks the sale as cancelled
func (s *Sale) Cancel() {
	s.SaleStatus = SaleStatusCancelled
}

// SaleItem represents an item in a sale
type SaleItem struct {
	ID               uint    `gorm:"primaryKey"`
	SaleID           uint    `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	ProductVariantID uint    `gorm:"not null;index"`
	Quantity         int     `gorm:"not null"`
	UnitPrice        float64 `gorm:"type:decimal(10,2);not null"`
	DiscountAmount   float64 `gorm:"type:decimal(10,2);default:0"`
	SubTotal         float64 `gorm:"type:decimal(10,2);not null"`
	TotalAmount      float64 `gorm:"type:decimal(10,2);not null"`
	CreatedAt        time.Time
}

func (SaleItem) TableName() string {
	return "sale_items"
}

// IsValid validates the sale item
func (si *SaleItem) IsValid() bool {
	return si.SaleID > 0 && si.ProductVariantID > 0 && si.Quantity > 0 && si.UnitPrice >= 0
}

// CalculateTotals calculates the totals for this item
func (si *SaleItem) CalculateTotals() {
	si.SubTotal = float64(si.Quantity) * si.UnitPrice
	si.TotalAmount = si.SubTotal - si.DiscountAmount
}

// PaymentMethod represents the method of payment
type PaymentMethod string

const (
	PaymentMethodCash   PaymentMethod = "cash"
	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodOther  PaymentMethod = "other"
)

func (pm PaymentMethod) IsValid() bool {
	switch pm {
	case PaymentMethodCash, PaymentMethodCard, PaymentMethodOther:
		return true
	}
	return false
}

// PaymentTransactionStatus represents the status of a payment transaction
type PaymentTransactionStatus string

const (
	PaymentTransactionStatusPending   PaymentTransactionStatus = "pending"
	PaymentTransactionStatusCompleted PaymentTransactionStatus = "completed"
	PaymentTransactionStatusFailed    PaymentTransactionStatus = "failed"
	PaymentTransactionStatusRefunded  PaymentTransactionStatus = "refunded"
)

func (pts PaymentTransactionStatus) IsValid() bool {
	switch pts {
	case PaymentTransactionStatusPending, PaymentTransactionStatusCompleted, 
		PaymentTransactionStatusFailed, PaymentTransactionStatusRefunded:
		return true
	}
	return false
}

// Payment represents a payment transaction for a sale
type Payment struct {
	ID            uint                     `gorm:"primaryKey"`
	SaleID        uint                     `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	PaymentMethod PaymentMethod            `gorm:"type:varchar(50);not null"`
	Amount        float64                  `gorm:"type:decimal(10,2);not null"`
	PaymentStatus PaymentTransactionStatus `gorm:"type:varchar(50);not null;default:'pending'"`
	Reference     string                   `gorm:"type:varchar(255)"` // Card transaction ref, check number, etc.
	Notes         string                   `gorm:"type:text"`
	CreatedAt     time.Time                `gorm:"index"`
}

func (Payment) TableName() string {
	return "payments"
}

// IsValid validates the payment
func (p *Payment) IsValid() bool {
	return p.SaleID > 0 && p.Amount > 0 && p.PaymentMethod.IsValid() && p.PaymentStatus.IsValid()
}

// MarkCompleted marks the payment as completed
func (p *Payment) MarkCompleted() {
	p.PaymentStatus = PaymentTransactionStatusCompleted
}

// CashDrawerStatus represents the status of a cash drawer
type CashDrawerStatus string

const (
	CashDrawerStatusOpen   CashDrawerStatus = "open"
	CashDrawerStatusClosed CashDrawerStatus = "closed"
)

func (cds CashDrawerStatus) IsValid() bool {
	switch cds {
	case CashDrawerStatusOpen, CashDrawerStatusClosed:
		return true
	}
	return false
}

// CashDrawer represents a cash drawer session
type CashDrawer struct {
	ID              uint             `gorm:"primaryKey"`
	CompanyID       *uint            `gorm:"index"`
	FranchiseID     *uint            `gorm:"index"`
	OpeningBalance  float64          `gorm:"type:decimal(10,2);not null"`
	ClosingBalance  *float64         `gorm:"type:decimal(10,2)"`
	ExpectedBalance *float64         `gorm:"type:decimal(10,2)"`
	Difference      *float64         `gorm:"type:decimal(10,2)"`
	Status          CashDrawerStatus `gorm:"type:varchar(50);not null;default:'open'"`
	OpenedByID      uint             `gorm:"not null;index"`
	ClosedByID      *uint            `gorm:"index"`
	OpenedAt        time.Time        `gorm:"not null;index"`
	ClosedAt        *time.Time
	Notes           string           `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Relationships
	Transactions []CashDrawerTransaction `gorm:"foreignKey:CashDrawerID"`
}

func (CashDrawer) TableName() string {
	return "cash_drawers"
}

// IsValid validates the cash drawer
func (cd *CashDrawer) IsValid() bool {
	companySet := cd.CompanyID != nil && *cd.CompanyID > 0
	franchiseSet := cd.FranchiseID != nil && *cd.FranchiseID > 0
	return cd.OpenedByID > 0 && ((companySet && !franchiseSet) || (!companySet && franchiseSet)) &&
		cd.Status.IsValid()
}

// IsOpen checks if the drawer is open
func (cd *CashDrawer) IsOpen() bool {
	return cd.Status == CashDrawerStatusOpen
}

// Close closes the cash drawer with the provided closing balance
func (cd *CashDrawer) Close(closingBalance float64, closedByID uint, expectedBalance float64) {
	cd.Status = CashDrawerStatusClosed
	cd.ClosingBalance = &closingBalance
	cd.ExpectedBalance = &expectedBalance
	difference := closingBalance - expectedBalance
	cd.Difference = &difference
	cd.ClosedByID = &closedByID
	now := time.Now()
	cd.ClosedAt = &now
}

// CashDrawerTransactionType represents the type of cash drawer transaction
type CashDrawerTransactionType string

const (
	CashDrawerTransactionTypeSale       CashDrawerTransactionType = "sale"
	CashDrawerTransactionTypeRefund     CashDrawerTransactionType = "refund"
	CashDrawerTransactionTypeAdjustment CashDrawerTransactionType = "adjustment"
)

func (cdtt CashDrawerTransactionType) IsValid() bool {
	switch cdtt {
	case CashDrawerTransactionTypeSale, CashDrawerTransactionTypeRefund, CashDrawerTransactionTypeAdjustment:
		return true
	}
	return false
}

// CashDrawerTransaction represents a transaction in the cash drawer
type CashDrawerTransaction struct {
	ID              uint                      `gorm:"primaryKey"`
	CashDrawerID    uint                      `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	TransactionType CashDrawerTransactionType `gorm:"type:varchar(50);not null"`
	Amount          float64                   `gorm:"type:decimal(10,2);not null"`
	SaleID          *uint                     `gorm:"index"` // Reference to sale if applicable
	Notes           string                    `gorm:"type:text"`
	CreatedAt       time.Time                 `gorm:"index"`
}

func (CashDrawerTransaction) TableName() string {
	return "cash_drawer_transactions"
}

// IsValid validates the cash drawer transaction
func (cdt *CashDrawerTransaction) IsValid() bool {
	return cdt.CashDrawerID > 0 && cdt.TransactionType.IsValid()
}

// RefundStatus represents the status of a refund
type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusCompleted RefundStatus = "completed"
	RefundStatusCancelled RefundStatus = "cancelled"
)

func (rs RefundStatus) IsValid() bool {
	switch rs {
	case RefundStatusPending, RefundStatusCompleted, RefundStatusCancelled:
		return true
	}
	return false
}

// Refund represents a refund transaction
type Refund struct {
	ID             uint          `gorm:"primaryKey"`
	OriginalSaleID uint          `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	RefundAmount   float64       `gorm:"type:decimal(10,2);not null"`
	Reason         string        `gorm:"type:text"`
	RefundMethod   PaymentMethod `gorm:"type:varchar(50);not null"`
	RefundStatus   RefundStatus  `gorm:"type:varchar(50);not null;default:'pending'"`
	ProcessedByID  uint          `gorm:"not null;index"`
	CreatedAt      time.Time     `gorm:"index"`
	UpdatedAt      time.Time

	// Relationships
	OriginalSale *Sale `gorm:"foreignKey:OriginalSaleID"`
}

func (Refund) TableName() string {
	return "refunds"
}

// IsValid validates the refund
func (r *Refund) IsValid() bool {
	return r.OriginalSaleID > 0 && r.RefundAmount > 0 && r.ProcessedByID > 0 &&
		r.RefundMethod.IsValid() && r.RefundStatus.IsValid()
}

// Complete marks the refund as completed
func (r *Refund) Complete() {
	r.RefundStatus = RefundStatusCompleted
}

// Cancel marks the refund as cancelled
func (r *Refund) Cancel() {
	r.RefundStatus = RefundStatusCancelled
}

