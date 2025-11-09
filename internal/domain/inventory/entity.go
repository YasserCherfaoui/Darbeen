package inventory

import "time"

type Inventory struct {
	ID               uint  `gorm:"primaryKey"`
	ProductVariantID uint  `gorm:"not null;index"`
	CompanyID        *uint `gorm:"index;constraint:OnDelete:CASCADE"` // Set if company inventory
	FranchiseID      *uint `gorm:"index;constraint:OnDelete:CASCADE"` // Set if franchise inventory
	Stock            int   `gorm:"default:0;not null"`
	ReservedStock    int   `gorm:"default:0;not null"`
	IsActive         bool  `gorm:"default:true"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Relationships (will be preloaded)
	ProductVariant interface{} `gorm:"-"` // Not stored, for service layer use
	Company        interface{} `gorm:"-"` // Not stored, for service layer use
	Franchise      interface{} `gorm:"-"` // Not stored, for service layer use
}

func (Inventory) TableName() string {
	return "inventories"
}

type MovementType string

const (
	MovementTypePurchase   MovementType = "purchase"
	MovementTypeSale       MovementType = "sale"
	MovementTypeAdjustment MovementType = "adjustment"
	MovementTypeTransfer   MovementType = "transfer"
	MovementTypeReserve    MovementType = "reserve"
	MovementTypeRelease    MovementType = "release"
	MovementTypeReturn     MovementType = "return"
)

func (mt MovementType) IsValid() bool {
	switch mt {
	case MovementTypePurchase, MovementTypeSale, MovementTypeAdjustment,
		MovementTypeTransfer, MovementTypeReserve, MovementTypeRelease, MovementTypeReturn:
		return true
	}
	return false
}

type InventoryMovement struct {
	ID            uint         `gorm:"primaryKey"`
	InventoryID   uint         `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	MovementType  MovementType `gorm:"type:varchar(50);not null;index"`
	Quantity      int          `gorm:"not null"`
	PreviousStock int          `gorm:"not null"`
	NewStock      int          `gorm:"not null"`
	ReferenceType *string      `gorm:"type:varchar(100)"` // e.g., 'order', 'transfer', 'adjustment'
	ReferenceID   *string      `gorm:"type:varchar(255)"` // Reference to external entity
	Notes         *string      `gorm:"type:text"`
	CreatedByID   uint         `gorm:"not null;index"`
	CreatedAt     time.Time    `gorm:"index"`
}

func (InventoryMovement) TableName() string {
	return "inventory_movements"
}

// Business methods for Inventory

// IsValid validates that either CompanyID or FranchiseID is set (but not both)
func (i *Inventory) IsValid() bool {
	companySet := i.CompanyID != nil && *i.CompanyID > 0
	franchiseSet := i.FranchiseID != nil && *i.FranchiseID > 0
	return i.ProductVariantID > 0 && ((companySet && !franchiseSet) || (!companySet && franchiseSet))
}

// BelongsToCompany checks if this inventory belongs to a specific company
func (i *Inventory) BelongsToCompany(companyID uint) bool {
	return i.CompanyID != nil && *i.CompanyID == companyID
}

// BelongsToFranchise checks if this inventory belongs to a specific franchise
func (i *Inventory) BelongsToFranchise(franchiseID uint) bool {
	return i.FranchiseID != nil && *i.FranchiseID == franchiseID
}

// GetAvailableStock returns the available stock (total stock minus reserved)
func (i *Inventory) GetAvailableStock() int {
	return i.Stock - i.ReservedStock
}

// CanFulfill checks if the available stock can fulfill the requested quantity
func (i *Inventory) CanFulfill(quantity int) bool {
	return i.GetAvailableStock() >= quantity && quantity > 0
}

// AddStock increases the stock level (for purchases, returns, etc.)
func (i *Inventory) AddStock(amount int) {
	if amount > 0 {
		i.Stock += amount
	}
}

// RemoveStock decreases the stock level if sufficient stock is available
func (i *Inventory) RemoveStock(amount int) bool {
	if amount > 0 && i.GetAvailableStock() >= amount {
		i.Stock -= amount
		return true
	}
	return false
}

// ReserveStock reserves stock for pending orders
func (i *Inventory) ReserveStock(amount int) bool {
	if amount > 0 && i.GetAvailableStock() >= amount {
		i.ReservedStock += amount
		return true
	}
	return false
}

// ReleaseStock releases previously reserved stock
func (i *Inventory) ReleaseStock(amount int) {
	if amount > 0 && i.ReservedStock >= amount {
		i.ReservedStock -= amount
	}
}

// FulfillReservation moves reserved stock to actual stock removal (for completed sales)
func (i *Inventory) FulfillReservation(amount int) bool {
	if amount > 0 && i.ReservedStock >= amount {
		i.ReservedStock -= amount
		return true
	}
	return false
}




