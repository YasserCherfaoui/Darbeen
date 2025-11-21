package warehousebill

import "time"

// BillType represents the type of warehouse bill
type BillType string

const (
	BillTypeExit  BillType = "exit"  // Warehouse → Franchise
	BillTypeEntry BillType = "entry"  // Franchise receives from warehouse
)

func (bt BillType) IsValid() bool {
	switch bt {
	case BillTypeExit, BillTypeEntry:
		return true
	}
	return false
}

// BillStatus represents the status of a warehouse bill
type BillStatus string

const (
	BillStatusDraft     BillStatus = "draft"
	BillStatusCompleted BillStatus = "completed"
	BillStatusCancelled BillStatus = "cancelled"
	BillStatusVerified  BillStatus = "verified"
)

func (bs BillStatus) IsValid() bool {
	switch bs {
	case BillStatusDraft, BillStatusCompleted, BillStatusCancelled, BillStatusVerified:
		return true
	}
	return false
}

// VerificationStatus represents the verification status of an entry bill
type VerificationStatus string

const (
	VerificationStatusPending            VerificationStatus = "pending"
	VerificationStatusVerified           VerificationStatus = "verified"
	VerificationStatusDiscrepanciesFound VerificationStatus = "discrepancies_found"
)

func (vs VerificationStatus) IsValid() bool {
	switch vs {
	case VerificationStatusPending, VerificationStatusVerified, VerificationStatusDiscrepanciesFound:
		return true
	}
	return false
}

// DiscrepancyType represents the type of discrepancy found during verification
type DiscrepancyType string

const (
	DiscrepancyTypeNone            DiscrepancyType = "none"
	DiscrepancyTypeMissing         DiscrepancyType = "missing"
	DiscrepancyTypeExtra           DiscrepancyType = "extra"
	DiscrepancyTypeQuantityMismatch DiscrepancyType = "quantity_mismatch"
)

func (dt DiscrepancyType) IsValid() bool {
	switch dt {
	case DiscrepancyTypeNone, DiscrepancyTypeMissing, DiscrepancyTypeExtra, DiscrepancyTypeQuantityMismatch:
		return true
	}
	return false
}

// WarehouseBill represents a bill for transferring inventory between warehouse and franchise
type WarehouseBill struct {
	ID                uint               `gorm:"primaryKey"`
	CompanyID         uint               `gorm:"not null;index"`
	FranchiseID        uint               `gorm:"not null;index"`
	BillNumber         string             `gorm:"uniqueIndex;not null"`
	BillType           BillType           `gorm:"type:varchar(50);not null;index"`
	RelatedBillID      *uint              `gorm:"index"` // For entry bills linking to exit bills
	Status             BillStatus          `gorm:"type:varchar(50);not null;default:'draft'"`
	VerificationStatus VerificationStatus `gorm:"type:varchar(50);default:'pending'"`
	TotalAmount        float64            `gorm:"type:decimal(10,2);not null"`
	Notes              string             `gorm:"type:text"`
	VerifiedByID       *uint              `gorm:"index"`
	VerifiedAt          *time.Time
	CreatedByID        uint               `gorm:"not null;index"`
	CreatedAt           time.Time          `gorm:"index"`
	UpdatedAt           time.Time

	// Relationships
	Items []WarehouseBillItem `gorm:"foreignKey:WarehouseBillID;constraint:OnDelete:CASCADE"`
}

func (WarehouseBill) TableName() string {
	return "warehouse_bills"
}

// IsValid validates the warehouse bill
func (wb *WarehouseBill) IsValid() bool {
	return wb.CompanyID > 0 && wb.FranchiseID > 0 &&
		wb.BillType.IsValid() && wb.Status.IsValid() &&
		wb.VerificationStatus.IsValid() && wb.CreatedByID > 0
}

// IsExitBill checks if this is an exit bill
func (wb *WarehouseBill) IsExitBill() bool {
	return wb.BillType == BillTypeExit
}

// IsEntryBill checks if this is an entry bill
func (wb *WarehouseBill) IsEntryBill() bool {
	return wb.BillType == BillTypeEntry
}

// CanBeVerified checks if the bill can be verified (must be entry bill in draft status)
func (wb *WarehouseBill) CanBeVerified() bool {
	return wb.IsEntryBill() && wb.Status == BillStatusDraft
}

// CanBeCompleted checks if the bill can be completed
func (wb *WarehouseBill) CanBeCompleted() bool {
	return wb.Status == BillStatusDraft || wb.Status == BillStatusVerified
}

// MarkVerified marks the bill as verified
func (wb *WarehouseBill) MarkVerified(verifiedByID uint) {
	wb.VerificationStatus = VerificationStatusVerified
	wb.Status = BillStatusVerified
	now := time.Now()
	wb.VerifiedByID = &verifiedByID
	wb.VerifiedAt = &now
}

// MarkDiscrepanciesFound marks the bill as having discrepancies
func (wb *WarehouseBill) MarkDiscrepanciesFound(verifiedByID uint) {
	wb.VerificationStatus = VerificationStatusDiscrepanciesFound
	wb.Status = BillStatusVerified
	now := time.Now()
	wb.VerifiedByID = &verifiedByID
	wb.VerifiedAt = &now
}

// Complete marks the bill as completed
func (wb *WarehouseBill) Complete() {
	wb.Status = BillStatusCompleted
}

// Cancel marks the bill as cancelled
func (wb *WarehouseBill) Cancel() {
	wb.Status = BillStatusCancelled
}

// WarehouseBillItem represents an item in a warehouse bill
type WarehouseBillItem struct {
	ID                uint            `gorm:"primaryKey"`
	WarehouseBillID   uint            `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	ProductVariantID  uint            `gorm:"not null;index"`
	ExpectedQuantity  int             `gorm:"default:0"` // For entry bills, copied from exit bill
	ReceivedQuantity  *int            // For entry bills, actual received (nullable)
	Quantity          int             `gorm:"not null"`  // For exit bills, quantity to send
	UnitPrice         float64         `gorm:"type:decimal(10,2);not null"`
	TotalAmount       float64         `gorm:"type:decimal(10,2);not null"`
	DiscrepancyType   DiscrepancyType `gorm:"type:varchar(50);default:'none'"`
	DiscrepancyNotes  string          `gorm:"type:text"`
	CreatedAt         time.Time

	// Relationships (for preloading)
	ProductVariant interface{} `gorm:"-"` // Not stored, for service layer use
}

func (WarehouseBillItem) TableName() string {
	return "warehouse_bill_items"
}

// IsValid validates the warehouse bill item
func (wbi *WarehouseBillItem) IsValid() bool {
	return wbi.WarehouseBillID > 0 && wbi.ProductVariantID > 0 &&
		wbi.Quantity > 0 && wbi.UnitPrice >= 0 &&
		wbi.DiscrepancyType.IsValid()
}

// CalculateTotal calculates the total amount for this item
func (wbi *WarehouseBillItem) CalculateTotal() {
	if wbi.ReceivedQuantity != nil {
		// For entry bills, use received quantity if available
		wbi.TotalAmount = float64(*wbi.ReceivedQuantity) * wbi.UnitPrice
	} else {
		// For exit bills or unverified entry bills, use expected/quantity
		qty := wbi.Quantity
		if wbi.ExpectedQuantity > 0 {
			qty = wbi.ExpectedQuantity
		}
		wbi.TotalAmount = float64(qty) * wbi.UnitPrice
	}
}

// SetDiscrepancy sets the discrepancy information for this item
func (wbi *WarehouseBillItem) SetDiscrepancy(discrepancyType DiscrepancyType, notes string) {
	wbi.DiscrepancyType = discrepancyType
	wbi.DiscrepancyNotes = notes
}

// HasDiscrepancy checks if this item has a discrepancy
func (wbi *WarehouseBillItem) HasDiscrepancy() bool {
	return wbi.DiscrepancyType != DiscrepancyTypeNone
}

