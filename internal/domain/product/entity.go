package product

import (
	"time"

	"gorm.io/datatypes"
)

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	CompanyID   uint   `gorm:"not null;index"`
	Name        string `gorm:"not null"`
	Description string
	SKU         string  `gorm:"not null"`
	BasePrice   float64 `gorm:"type:decimal(10,2)"`
	IsActive    bool    `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relationships
	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string {
	return "products"
}

type ProductVariant struct {
	ID         uint           `gorm:"primaryKey"`
	ProductID  uint           `gorm:"not null;index"`
	Name       string         `gorm:"not null"`
	SKU        string         `gorm:"not null"`
	Price      float64        `gorm:"type:decimal(10,2)"`
	Stock      int            `gorm:"default:0"`
	Attributes datatypes.JSON `gorm:"type:jsonb"`
	IsActive   bool           `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Relationships
	Product *Product `gorm:"foreignKey:ProductID"`
}

func (ProductVariant) TableName() string {
	return "product_variants"
}

// Business methods for Product
func (p *Product) IsValid() bool {
	return p.Name != "" && p.SKU != "" && p.CompanyID > 0
}

func (p *Product) GetEffectivePrice() float64 {
	if p.BasePrice > 0 {
		return p.BasePrice
	}
	return 0
}

// Business methods for ProductVariant
func (pv *ProductVariant) IsValid() bool {
	return pv.Name != "" && pv.SKU != "" && pv.ProductID > 0
}

func (pv *ProductVariant) GetEffectivePrice(basePrice float64) float64 {
	if pv.Price > 0 {
		return pv.Price
	}
	return basePrice
}

func (pv *ProductVariant) HasStock() bool {
	return pv.Stock > 0
}

func (pv *ProductVariant) UpdateStock(newStock int) {
	pv.Stock = newStock
}

func (pv *ProductVariant) AddStock(amount int) {
	pv.Stock += amount
}

func (pv *ProductVariant) RemoveStock(amount int) bool {
	if pv.Stock >= amount {
		pv.Stock -= amount
		return true
	}
	return false
}
