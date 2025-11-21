package otp

import "time"

type OTP struct {
	ID          uint      `gorm:"primaryKey"`
	Code        string    `gorm:"not null;index;size:6"` // 6-digit numeric code
	Email       string    `gorm:"not null;index"`
	UserID      uint      `gorm:"not null;index"`
	CompanyID   uint      `gorm:"not null;index"`
	FranchiseID *uint     `gorm:"index"` // Nullable
	Purpose     Purpose   `gorm:"type:varchar(50);not null;index"`
	Used        bool      `gorm:"default:false;index"`
	ExpiresAt   time.Time `gorm:"not null;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (OTP) TableName() string {
	return "otps"
}

type Purpose string

const (
	PurposeSetupAccount  Purpose = "setup_account"
	PurposePasswordReset Purpose = "password_reset"
)

func (p Purpose) String() string {
	return string(p)
}

func (p Purpose) IsValid() bool {
	switch p {
	case PurposeSetupAccount, PurposePasswordReset:
		return true
	}
	return false
}

