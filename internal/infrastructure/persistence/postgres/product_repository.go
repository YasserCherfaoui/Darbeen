package postgres

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/product"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) product.Repository {
	return &productRepository{db: db}
}

// Product operations
func (r *productRepository) CreateProduct(p *product.Product) error {
	return r.db.Create(p).Error
}

func (r *productRepository) FindProductByID(id uint) (*product.Product, error) {
	var p product.Product
	err := r.db.Preload("Variants").Where("id = ?", id).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) FindProductByIDAndCompany(id, companyID uint) (*product.Product, error) {
	var p product.Product
	err := r.db.Preload("Variants").Where("id = ? AND company_id = ?", id, companyID).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) FindProductBySKUAndCompany(sku string, companyID uint) (*product.Product, error) {
	var p product.Product
	err := r.db.Where("sku = ? AND company_id = ?", sku, companyID).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) FindProductsByCompanyID(companyID uint, page, limit int) ([]*product.Product, int64, error) {
	var products []*product.Product
	var total int64

	// Count total
	err := r.db.Model(&product.Product{}).Where("company_id = ? AND is_active = ?", companyID, true).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch products with variants
	err = r.db.Preload("Variants", "is_active = ?", true).
		Where("company_id = ? AND is_active = ?", companyID, true).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) UpdateProduct(p *product.Product) error {
	return r.db.Save(p).Error
}

func (r *productRepository) SoftDeleteProduct(id uint) error {
	return r.db.Model(&product.Product{}).Where("id = ?", id).Update("is_active", false).Error
}

// Product variant operations
func (r *productRepository) CreateProductVariant(v *product.ProductVariant) error {
	return r.db.Create(v).Error
}

func (r *productRepository) FindProductVariantByID(id uint) (*product.ProductVariant, error) {
	var v product.ProductVariant
	err := r.db.Where("id = ?", id).First(&v).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product variant not found")
		}
		return nil, err
	}
	return &v, nil
}

func (r *productRepository) FindProductVariantByIDAndProduct(id, productID uint) (*product.ProductVariant, error) {
	var v product.ProductVariant
	err := r.db.Where("id = ? AND product_id = ?", id, productID).First(&v).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product variant not found")
		}
		return nil, err
	}
	return &v, nil
}

func (r *productRepository) FindProductVariantBySKUAndProduct(sku string, productID uint) (*product.ProductVariant, error) {
	var v product.ProductVariant
	err := r.db.Where("sku = ? AND product_id = ?", sku, productID).First(&v).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product variant not found")
		}
		return nil, err
	}
	return &v, nil
}

func (r *productRepository) FindProductVariantsByProductID(productID uint) ([]*product.ProductVariant, error) {
	var variants []*product.ProductVariant
	err := r.db.Where("product_id = ? AND is_active = ?", productID, true).Find(&variants).Error
	return variants, err
}

func (r *productRepository) UpdateProductVariant(v *product.ProductVariant) error {
	return r.db.Save(v).Error
}

func (r *productRepository) SoftDeleteProductVariant(id uint) error {
	return r.db.Model(&product.ProductVariant{}).Where("id = ?", id).Update("is_active", false).Error
}

// Stock operations
func (r *productRepository) UpdateVariantStock(variantID uint, newStock int) error {
	return r.db.Model(&product.ProductVariant{}).Where("id = ?", variantID).Update("stock", newStock).Error
}

func (r *productRepository) AddVariantStock(variantID uint, amount int) error {
	return r.db.Model(&product.ProductVariant{}).Where("id = ?", variantID).Update("stock", gorm.Expr("stock + ?", amount)).Error
}

func (r *productRepository) RemoveVariantStock(variantID uint, amount int) error {
	return r.db.Model(&product.ProductVariant{}).Where("id = ? AND stock >= ?", variantID, amount).Update("stock", gorm.Expr("stock - ?", amount)).Error
}
