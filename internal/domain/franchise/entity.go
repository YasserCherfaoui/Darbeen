package franchise

import "time"

type Franchise struct {
	ID              uint   `gorm:"primaryKey"`
	ParentCompanyID uint   `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Name            string `gorm:"not null"`
	Code            string `gorm:"uniqueIndex;not null"`
	Description     string
	Address         string
	Phone           string
	Email           string
	IsActive        bool `gorm:"default:true"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Relationships (will be preloaded)
	ParentCompany interface{} `gorm:"-"` // Not stored, for service layer use
}

func (Franchise) TableName() string {
	return "franchises"
}

type FranchisePricing struct {
	ID               uint     `gorm:"primaryKey"`
	FranchiseID      uint     `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	ProductVariantID uint     `gorm:"not null;index"`
	RetailPrice      *float64 `gorm:"type:decimal(10,2)"` // Nullable - overrides variant retail price
	WholesalePrice   *float64 `gorm:"type:decimal(10,2)"` // Nullable - overrides variant wholesale price
	IsActive         bool     `gorm:"default:true"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Relationships (will be preloaded)
	Franchise      interface{} `gorm:"-"` // Not stored, for service layer use
	ProductVariant interface{} `gorm:"-"` // Not stored, for service layer use
}

func (FranchisePricing) TableName() string {
	return "franchise_pricing"
}

// Business methods for Franchise

// IsValid validates the franchise entity
func (f *Franchise) IsValid() bool {
	return f.Name != "" && f.Code != "" && f.ParentCompanyID > 0
}

// BelongsToCompany checks if this franchise belongs to a specific company
func (f *Franchise) BelongsToCompany(companyID uint) bool {
	return f.ParentCompanyID == companyID
}

// SetRetailPriceOverride sets the retail price override for this franchise
func (fp *FranchisePricing) SetRetailPriceOverride(price float64) {
	if price > 0 {
		fp.RetailPrice = &price
	} else {
		fp.RetailPrice = nil
	}
}

// SetWholesalePriceOverride sets the wholesale price override for this franchise
func (fp *FranchisePricing) SetWholesalePriceOverride(price float64) {
	if price > 0 {
		fp.WholesalePrice = &price
	} else {
		fp.WholesalePrice = nil
	}
}

// HasRetailOverride checks if there's a retail price override
func (fp *FranchisePricing) HasRetailOverride() bool {
	return fp.RetailPrice != nil && *fp.RetailPrice > 0
}

// HasWholesaleOverride checks if there's a wholesale price override
func (fp *FranchisePricing) HasWholesaleOverride() bool {
	return fp.WholesalePrice != nil && *fp.WholesalePrice > 0
}

// GetEffectiveRetailPrice returns the effective retail price (override or nil if none)
func (fp *FranchisePricing) GetEffectiveRetailPrice() *float64 {
	if fp.RetailPrice != nil && *fp.RetailPrice > 0 {
		return fp.RetailPrice
	}
	return nil
}

// GetEffectiveWholesalePrice returns the effective wholesale price (override or nil if none)
func (fp *FranchisePricing) GetEffectiveWholesalePrice() *float64 {
	if fp.WholesalePrice != nil && *fp.WholesalePrice > 0 {
		return fp.WholesalePrice
	}
	return nil
}




