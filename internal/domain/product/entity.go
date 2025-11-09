package product

import (
	"time"

	"gorm.io/datatypes"
)

type Product struct {
	ID                 uint   `gorm:"primaryKey"`
	CompanyID          uint   `gorm:"not null;index"`
	Name               string `gorm:"not null"`
	Description        string
	SKU                string  `gorm:"not null"`
	BaseRetailPrice    float64 `gorm:"type:decimal(10,2);not null"`
	BaseWholesalePrice float64 `gorm:"type:decimal(10,2);not null"`
	SupplierID         *uint   `gorm:"index"`                 // Nullable - product may not have a supplier
	SupplierCost       *float64 `gorm:"type:decimal(10,2)"`    // Nullable - cost from supplier
	IsActive           bool    `gorm:"default:true"`
	CreatedAt          time.Time
	UpdatedAt          time.Time

	// Relationships
	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string {
	return "products"
}

type ProductVariant struct {
	ID               uint           `gorm:"primaryKey"`
	ProductID        uint           `gorm:"not null;index"`
	Name             string         `gorm:"not null"`
	SKU              string         `gorm:"not null"`
	RetailPrice      *float64       `gorm:"type:decimal(10,2)"` // Nullable - inherits from parent if nil
	WholesalePrice   *float64       `gorm:"type:decimal(10,2)"` // Nullable - inherits from parent if nil
	UseParentPricing bool           `gorm:"default:false"`
	Attributes       datatypes.JSON `gorm:"type:jsonb"`
	IsActive         bool           `gorm:"default:true"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Relationships
	Product *Product `gorm:"foreignKey:ProductID"`
}

func (ProductVariant) TableName() string {
	return "product_variants"
}

// Business methods for Product
func (p *Product) IsValid() bool {
	return p.Name != "" && p.SKU != "" && p.CompanyID > 0 && p.BaseRetailPrice >= 0 && p.BaseWholesalePrice >= 0
}

func (p *Product) GetBaseRetailPrice() float64 {
	return p.BaseRetailPrice
}

func (p *Product) GetBaseWholesalePrice() float64 {
	return p.BaseWholesalePrice
}

// Business methods for ProductVariant
func (pv *ProductVariant) IsValid() bool {
	return pv.Name != "" && pv.SKU != "" && pv.ProductID > 0
}

// GetEffectiveRetailPrice returns the effective retail price (variant-specific or inherited from product)
func (pv *ProductVariant) GetEffectiveRetailPrice(baseRetailPrice float64) float64 {
	if pv.UseParentPricing {
		return baseRetailPrice
	}
	if pv.RetailPrice != nil && *pv.RetailPrice > 0 {
		return *pv.RetailPrice
	}
	return baseRetailPrice
}

// GetEffectiveWholesalePrice returns the effective wholesale price (variant-specific or inherited from product)
func (pv *ProductVariant) GetEffectiveWholesalePrice(baseWholesalePrice float64) float64 {
	if pv.UseParentPricing {
		return baseWholesalePrice
	}
	if pv.WholesalePrice != nil && *pv.WholesalePrice > 0 {
		return *pv.WholesalePrice
	}
	return baseWholesalePrice
}

// SetRetailPrice sets the retail price for this variant
func (pv *ProductVariant) SetRetailPrice(price float64) {
	if price > 0 {
		pv.RetailPrice = &price
		pv.UseParentPricing = false
	}
}

// SetWholesalePrice sets the wholesale price for this variant
func (pv *ProductVariant) SetWholesalePrice(price float64) {
	if price > 0 {
		pv.WholesalePrice = &price
		pv.UseParentPricing = false
	}
}

// UseParentPricing marks the variant to inherit pricing from parent product
func (pv *ProductVariant) MarkUseParentPricing() {
	pv.UseParentPricing = true
	pv.RetailPrice = nil
	pv.WholesalePrice = nil
}
