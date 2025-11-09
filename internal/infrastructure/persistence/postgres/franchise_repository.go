package postgres

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/franchise"
	"gorm.io/gorm"
)

type franchiseRepository struct {
	db *gorm.DB
}

func NewFranchiseRepository(db *gorm.DB) franchise.Repository {
	return &franchiseRepository{db: db}
}

// Franchise CRUD

func (r *franchiseRepository) Create(f *franchise.Franchise) error {
	return r.db.Create(f).Error
}

func (r *franchiseRepository) FindByID(id uint) (*franchise.Franchise, error) {
	var f franchise.Franchise
	err := r.db.Where("id = ?", id).First(&f).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("franchise not found")
		}
		return nil, err
	}
	return &f, nil
}

func (r *franchiseRepository) FindByCode(code string) (*franchise.Franchise, error) {
	var f franchise.Franchise
	err := r.db.Where("code = ?", code).First(&f).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("franchise not found")
		}
		return nil, err
	}
	return &f, nil
}

func (r *franchiseRepository) FindByParentCompanyID(companyID uint) ([]*franchise.Franchise, error) {
	var franchises []*franchise.Franchise
	err := r.db.Where("parent_company_id = ? AND is_active = ?", companyID, true).Find(&franchises).Error
	return franchises, err
}

func (r *franchiseRepository) Update(f *franchise.Franchise) error {
	return r.db.Save(f).Error
}

// Franchise Pricing

func (r *franchiseRepository) CreatePricing(pricing *franchise.FranchisePricing) error {
	return r.db.Create(pricing).Error
}

func (r *franchiseRepository) FindPricing(franchiseID, variantID uint) (*franchise.FranchisePricing, error) {
	var pricing franchise.FranchisePricing
	err := r.db.Where("franchise_id = ? AND product_variant_id = ? AND is_active = ?", franchiseID, variantID, true).First(&pricing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("franchise pricing not found")
		}
		return nil, err
	}
	return &pricing, nil
}

func (r *franchiseRepository) FindAllPricingByFranchise(franchiseID uint) ([]*franchise.FranchisePricing, error) {
	var pricings []*franchise.FranchisePricing
	err := r.db.Where("franchise_id = ? AND is_active = ?", franchiseID, true).Find(&pricings).Error
	return pricings, err
}

func (r *franchiseRepository) UpdatePricing(pricing *franchise.FranchisePricing) error {
	return r.db.Save(pricing).Error
}

func (r *franchiseRepository) DeletePricing(franchiseID, variantID uint) error {
	return r.db.Model(&franchise.FranchisePricing{}).
		Where("franchise_id = ? AND product_variant_id = ?", franchiseID, variantID).
		Update("is_active", false).Error
}




