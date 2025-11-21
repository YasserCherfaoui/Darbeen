package postgres

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/otp"
	"gorm.io/gorm"
)

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) otp.Repository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(otpRecord *otp.OTP) error {
	return r.db.Create(otpRecord).Error
}

func (r *otpRepository) FindByCodeAndEmail(code, email string) (*otp.OTP, error) {
	var otpRecord otp.OTP
	err := r.db.Where("code = ? AND email = ? AND used = ?", code, email, false).First(&otpRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("otp not found")
		}
		return nil, err
	}
	return &otpRecord, nil
}

func (r *otpRepository) FindByCode(code string) (*otp.OTP, error) {
	var otpRecord otp.OTP
	err := r.db.Where("code = ? AND used = ?", code, false).First(&otpRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("otp not found")
		}
		return nil, err
	}
	return &otpRecord, nil
}

func (r *otpRepository) MarkAsUsed(id uint) error {
	return r.db.Model(&otp.OTP{}).Where("id = ?", id).Update("used", true).Error
}

func (r *otpRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ? OR used = ?", "NOW()", true).Delete(&otp.OTP{}).Error
}

