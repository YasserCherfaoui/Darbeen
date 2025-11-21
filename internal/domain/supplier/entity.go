package supplier

import (
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/product"
)

type Supplier struct {
	ID            uint   `gorm:"primaryKey"`
	CompanyID     uint   `gorm:"not null;index"`
	Name          string `gorm:"not null"`
	ContactPerson string
	Email         string
	Phone         string
	Address       string
	IsActive      bool `gorm:"default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (Supplier) TableName() string {
	return "suppliers"
}

// Business methods for Supplier
func (s *Supplier) IsValid() bool {
	return s.Name != "" && s.CompanyID > 0
}

// BillStatus represents the status of a supplier bill
type BillStatus string

const (
	BillStatusDraft     BillStatus = "draft"
	BillStatusCompleted BillStatus = "completed"
	BillStatusCancelled BillStatus = "cancelled"
)

func (bs BillStatus) IsValid() bool {
	switch bs {
	case BillStatusDraft, BillStatusCompleted, BillStatusCancelled:
		return true
	}
	return false
}

// PaymentStatus represents the payment status of a bill
type PaymentStatus string

const (
	PaymentStatusUnpaid       PaymentStatus = "unpaid"
	PaymentStatusPartiallyPaid PaymentStatus = "partially_paid"
	PaymentStatusPaid         PaymentStatus = "paid"
)

func (ps PaymentStatus) IsValid() bool {
	switch ps {
	case PaymentStatusUnpaid, PaymentStatusPartiallyPaid, PaymentStatusPaid:
		return true
	}
	return false
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

// SupplierBill represents a bill from a supplier
type SupplierBill struct {
	ID            uint          `gorm:"primaryKey"`
	CompanyID     uint          `gorm:"not null;index"`
	SupplierID    uint          `gorm:"not null;index"`
	BillNumber    string        `gorm:"uniqueIndex;not null"`
	TotalAmount   float64       `gorm:"type:decimal(10,2);not null"`
	PaidAmount    float64       `gorm:"type:decimal(10,2);default:0;not null"`
	PendingAmount float64       `gorm:"type:decimal(10,2);not null"` // Calculated: TotalAmount - PaidAmount
	PaymentStatus PaymentStatus `gorm:"type:varchar(50);not null;default:'unpaid'"`
	BillStatus    BillStatus    `gorm:"type:varchar(50);not null;default:'draft'"`
	Notes         string        `gorm:"type:text"`
	CreatedByID   uint          `gorm:"not null;index"`
	CreatedAt     time.Time     `gorm:"index"`
	UpdatedAt     time.Time

	// Relationships
	Supplier *Supplier              `gorm:"foreignKey:SupplierID"`
	Items    []SupplierBillItem     `gorm:"foreignKey:SupplierBillID;constraint:OnDelete:CASCADE"`
}

func (SupplierBill) TableName() string {
	return "supplier_bills"
}

// IsValid validates the supplier bill
// Note: BillNumber is optional during creation (will be auto-generated)
func (sb *SupplierBill) IsValid() bool {
	return sb.CompanyID > 0 && sb.SupplierID > 0 &&
		sb.TotalAmount >= 0 && sb.PaidAmount >= 0 && sb.CreatedByID > 0 &&
		sb.PaymentStatus.IsValid() && sb.BillStatus.IsValid()
}

// CalculatePendingAmount calculates and updates the pending amount
func (sb *SupplierBill) CalculatePendingAmount() {
	sb.PendingAmount = sb.TotalAmount - sb.PaidAmount
	if sb.PendingAmount < 0 {
		sb.PendingAmount = 0
	}
}

// UpdatePaymentStatus updates the payment status based on paid vs total amount
func (sb *SupplierBill) UpdatePaymentStatus() {
	sb.CalculatePendingAmount()
	if sb.PaidAmount == 0 {
		sb.PaymentStatus = PaymentStatusUnpaid
	} else if sb.PendingAmount <= 0 {
		sb.PaymentStatus = PaymentStatusPaid
	} else {
		sb.PaymentStatus = PaymentStatusPartiallyPaid
	}
}

// AddPayment adds to the paid amount
func (sb *SupplierBill) AddPayment(amount float64) {
	if amount > 0 {
		sb.PaidAmount += amount
		sb.UpdatePaymentStatus()
	}
}

// Complete marks the bill as completed
func (sb *SupplierBill) Complete() {
	sb.BillStatus = BillStatusCompleted
}

// Cancel marks the bill as cancelled
func (sb *SupplierBill) Cancel() {
	sb.BillStatus = BillStatusCancelled
}

// SupplierBillItem represents an item in a supplier bill
type SupplierBillItem struct {
	ID               uint    `gorm:"primaryKey"`
	SupplierBillID   uint    `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	ProductVariantID uint    `gorm:"not null;index"`
	Quantity         int     `gorm:"not null"`
	UnitCost         float64 `gorm:"type:decimal(10,2);not null"`
	TotalCost        float64 `gorm:"type:decimal(10,2);not null"`
	CreatedAt        time.Time

	// Relationships (for preloading)
	ProductVariant *product.ProductVariant `gorm:"foreignKey:ProductVariantID"`
}

func (SupplierBillItem) TableName() string {
	return "supplier_bill_items"
}

// IsValid validates the bill item
func (sbi *SupplierBillItem) IsValid() bool {
	return sbi.SupplierBillID > 0 && sbi.ProductVariantID > 0 &&
		sbi.Quantity > 0 && sbi.UnitCost >= 0
}

// CalculateTotal calculates the total cost for this item
func (sbi *SupplierBillItem) CalculateTotal() {
	sbi.TotalCost = float64(sbi.Quantity) * sbi.UnitCost
}

// SupplierPayment represents a payment made to a supplier
type SupplierPayment struct {
	ID            uint                     `gorm:"primaryKey"`
	SupplierID    uint                     `gorm:"not null;index"`
	CompanyID     uint                     `gorm:"not null;index"`
	Amount        float64                  `gorm:"type:decimal(10,2);not null"`
	PaymentMethod PaymentMethod            `gorm:"type:varchar(50);not null"`
	PaymentStatus PaymentTransactionStatus `gorm:"type:varchar(50);not null;default:'pending'"`
	Reference     string                   `gorm:"type:varchar(255)"`
	Notes         string                   `gorm:"type:text"`
	CreatedByID   uint                     `gorm:"not null;index"`
	CreatedAt     time.Time                `gorm:"index"`

	// Relationships
	Distributions []SupplierPaymentDistribution `gorm:"foreignKey:SupplierPaymentID;constraint:OnDelete:CASCADE"`
}

func (SupplierPayment) TableName() string {
	return "supplier_payments"
}

// IsValid validates the supplier payment
func (sp *SupplierPayment) IsValid() bool {
	return sp.SupplierID > 0 && sp.CompanyID > 0 && sp.Amount > 0 &&
		sp.PaymentMethod.IsValid() && sp.PaymentStatus.IsValid() && sp.CreatedByID > 0
}

// MarkCompleted marks the payment as completed
func (sp *SupplierPayment) MarkCompleted() {
	sp.PaymentStatus = PaymentTransactionStatusCompleted
}

// SupplierPaymentDistribution represents how a payment is distributed across bills
type SupplierPaymentDistribution struct {
	ID                uint      `gorm:"primaryKey"`
	SupplierPaymentID uint      `gorm:"not null;index"`
	SupplierBillID    uint      `gorm:"not null;index"`
	Amount            float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt         time.Time

	// Relationships for foreign key constraints
	SupplierPayment *SupplierPayment `gorm:"foreignKey:SupplierPaymentID;references:ID;constraint:OnDelete:CASCADE"`
	SupplierBill    *SupplierBill    `gorm:"foreignKey:SupplierBillID;references:ID;constraint:OnDelete:RESTRICT"`
}

func (SupplierPaymentDistribution) TableName() string {
	return "supplier_payment_distributions"
}

// IsValid validates the payment distribution
func (spd *SupplierPaymentDistribution) IsValid() bool {
	return spd.SupplierPaymentID > 0 && spd.SupplierBillID > 0 && spd.Amount > 0
}

