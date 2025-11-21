package otp

import (
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/otp"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/YasserCherfaoui/darween/pkg/token"
)

type Service struct {
	otpRepo otp.Repository
}

func NewService(otpRepo otp.Repository) *Service {
	return &Service{
		otpRepo: otpRepo,
	}
}

// GenerateOTP generates a new OTP and stores it in the database
func (s *Service) GenerateOTP(email string, userID, companyID uint, franchiseID *uint, purpose otp.Purpose, expiryDuration time.Duration) (string, error) {
	// Generate 6-digit numeric OTP
	code, err := token.GenerateOTP()
	if err != nil {
		return "", errors.NewInternalError("failed to generate OTP", err)
	}

	// Create OTP record
	otpRecord := &otp.OTP{
		Code:        code,
		Email:       email,
		UserID:      userID,
		CompanyID:   companyID,
		FranchiseID: franchiseID,
		Purpose:     purpose,
		Used:        false,
		ExpiresAt:   time.Now().Add(expiryDuration),
	}

	if err := s.otpRepo.Create(otpRecord); err != nil {
		return "", errors.NewInternalError("failed to store OTP", err)
	}

	return code, nil
}

// ValidateOTP validates an OTP code for a given email
func (s *Service) ValidateOTP(code, email string) (*otp.OTP, error) {
	// Find OTP by code and email
	otpRecord, err := s.otpRepo.FindByCodeAndEmail(code, email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid OTP code")
	}

	// Check if OTP is expired
	if time.Now().After(otpRecord.ExpiresAt) {
		return nil, errors.NewUnauthorizedError("OTP code has expired")
	}

	// Check if OTP is already used
	if otpRecord.Used {
		return nil, errors.NewUnauthorizedError("OTP code has already been used")
	}

	return otpRecord, nil
}

// MarkAsUsed marks an OTP as used
func (s *Service) MarkAsUsed(otpID uint) error {
	return s.otpRepo.MarkAsUsed(otpID)
}

// CleanupExpiredOTPs removes expired and used OTPs from the database
func (s *Service) CleanupExpiredOTPs() error {
	return s.otpRepo.DeleteExpired()
}

