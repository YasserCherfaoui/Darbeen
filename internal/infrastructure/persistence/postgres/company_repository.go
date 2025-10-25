package postgres

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/company"
	"gorm.io/gorm"
)

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) company.Repository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(c *company.Company) error {
	return r.db.Create(c).Error
}

func (r *companyRepository) FindByID(id uint) (*company.Company, error) {
	var c company.Company
	err := r.db.Where("id = ?", id).First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("company not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *companyRepository) FindByCode(code string) (*company.Company, error) {
	var c company.Company
	err := r.db.Where("code = ?", code).First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("company not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *companyRepository) Update(c *company.Company) error {
	return r.db.Save(c).Error
}

func (r *companyRepository) FindByUserID(userID uint) ([]*company.Company, error) {
	var companies []*company.Company
	err := r.db.
		Joins("JOIN user_company_roles ON user_company_roles.company_id = companies.id").
		Where("user_company_roles.user_id = ? AND user_company_roles.is_active = ?", userID, true).
		Find(&companies).Error
	return companies, err
}
