package postgres

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/smtpconfig"
	"gorm.io/gorm"
)

type smtpConfigRepository struct {
	db *gorm.DB
}

func NewSMTPConfigRepository(db *gorm.DB) smtpconfig.Repository {
	return &smtpConfigRepository{db: db}
}

func (r *smtpConfigRepository) Create(config *smtpconfig.SMTPConfig) error {
	return r.db.Create(config).Error
}

func (r *smtpConfigRepository) FindByID(id uint) (*smtpconfig.SMTPConfig, error) {
	var config smtpconfig.SMTPConfig
	err := r.db.Where("id = ?", id).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("smtp config not found")
		}
		return nil, err
	}
	return &config, nil
}

func (r *smtpConfigRepository) FindByIDAndCompany(id, companyID uint) (*smtpconfig.SMTPConfig, error) {
	var config smtpconfig.SMTPConfig
	err := r.db.Where("id = ? AND company_id = ?", id, companyID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("smtp config not found")
		}
		return nil, err
	}
	return &config, nil
}

func (r *smtpConfigRepository) FindByCompanyID(companyID uint) ([]*smtpconfig.SMTPConfig, error) {
	var configs []*smtpconfig.SMTPConfig
	err := r.db.Where("company_id = ?", companyID).Find(&configs).Error
	return configs, err
}

func (r *smtpConfigRepository) FindDefaultByCompanyID(companyID uint) (*smtpconfig.SMTPConfig, error) {
	var config smtpconfig.SMTPConfig
	err := r.db.Where("company_id = ? AND is_default = ? AND is_active = ?", companyID, true, true).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("default smtp config not found")
		}
		return nil, err
	}
	return &config, nil
}

func (r *smtpConfigRepository) Update(config *smtpconfig.SMTPConfig) error {
	return r.db.Save(config).Error
}

func (r *smtpConfigRepository) Delete(id uint) error {
	return r.db.Delete(&smtpconfig.SMTPConfig{}, id).Error
}

func (r *smtpConfigRepository) SetAsDefault(id, companyID uint) error {
	// Start transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Unset all defaults for this company
	if err := tx.Model(&smtpconfig.SMTPConfig{}).
		Where("company_id = ?", companyID).
		Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the specified config as default
	if err := tx.Model(&smtpconfig.SMTPConfig{}).
		Where("id = ? AND company_id = ?", id, companyID).
		Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

