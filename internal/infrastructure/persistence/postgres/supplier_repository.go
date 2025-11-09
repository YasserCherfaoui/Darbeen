package postgres

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/product"
	"github.com/YasserCherfaoui/darween/internal/domain/supplier"
	"gorm.io/gorm"
)

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) supplier.Repository {
	return &supplierRepository{db: db}
}

// Supplier operations
func (r *supplierRepository) CreateSupplier(s *supplier.Supplier) error {
	return r.db.Create(s).Error
}

func (r *supplierRepository) FindSupplierByID(id uint) (*supplier.Supplier, error) {
	var s supplier.Supplier
	err := r.db.Where("id = ?", id).First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *supplierRepository) FindSupplierByIDAndCompany(id, companyID uint) (*supplier.Supplier, error) {
	var s supplier.Supplier
	err := r.db.Where("id = ? AND company_id = ?", id, companyID).First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *supplierRepository) FindSuppliersByCompanyID(companyID uint, page, limit int) ([]*supplier.Supplier, int64, error) {
	var suppliers []*supplier.Supplier
	var total int64

	// Count total
	err := r.db.Model(&supplier.Supplier{}).Where("company_id = ? AND is_active = ?", companyID, true).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch suppliers
	err = r.db.Where("company_id = ? AND is_active = ?", companyID, true).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&suppliers).Error

	return suppliers, total, err
}

func (r *supplierRepository) UpdateSupplier(s *supplier.Supplier) error {
	return r.db.Save(s).Error
}

func (r *supplierRepository) SoftDeleteSupplier(id uint) error {
	return r.db.Model(&supplier.Supplier{}).Where("id = ?", id).Update("is_active", false).Error
}

// Product-Supplier relationship
func (r *supplierRepository) FindProductsBySupplier(supplierID uint) ([]*product.Product, error) {
	var products []*product.Product
	err := r.db.Where("supplier_id = ? AND is_active = ?", supplierID, true).Find(&products).Error
	return products, err
}

