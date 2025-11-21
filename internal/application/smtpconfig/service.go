package smtpconfig

import (
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/smtpconfig"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/pkg/encryption"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	smtpRepo smtpconfig.Repository
	userRepo user.Repository
}

func NewService(smtpRepo smtpconfig.Repository, userRepo user.Repository) *Service {
	return &Service{
		smtpRepo: smtpRepo,
		userRepo: userRepo,
	}
}

func (s *Service) CreateSMTPConfig(userID, companyID uint, req *CreateSMTPConfigRequest) (*SMTPConfigResponse, error) {
	// Check user role in company
	role, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	if role.Role != user.RoleOwner && role.Role != user.RoleAdmin {
		return nil, errors.NewForbiddenError("only owners and admins can manage SMTP configs")
	}

	// Validate security type
	securityType := smtpconfig.SecurityType(req.Security)
	if !securityType.IsValid() {
		return nil, errors.NewValidationError("invalid security type")
	}

	// Validate port range
	if req.Port < 1 || req.Port > 65535 {
		return nil, errors.NewValidationError("port must be between 1 and 65535")
	}

	// Set default rate limit if not provided
	rateLimit := req.RateLimit
	if rateLimit == 0 {
		rateLimit = 80 // Default 80 emails per hour
	}

	// Encrypt password (using AES so it can be decrypted for SMTP)
	encryptedPassword, err := encryption.Encrypt(req.Password)
	if err != nil {
		return nil, errors.NewInternalError("failed to encrypt password", err)
	}

	// Set default IsActive if not provided
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Create SMTP config
	newConfig := &smtpconfig.SMTPConfig{
		CompanyID: companyID,
		Host:      req.Host,
		User:      req.User,
		Password:  encryptedPassword,
		Port:      req.Port,
		FromName:  req.FromName,
		Security:  securityType,
		RateLimit: rateLimit,
		IsActive:  isActive,
		IsDefault: false,
	}

	if err := s.smtpRepo.Create(newConfig); err != nil {
		return nil, errors.NewInternalError("failed to create SMTP config", err)
	}

	return s.toResponse(newConfig), nil
}

func (s *Service) GetSMTPConfigs(userID, companyID uint) (*ListSMTPConfigsResponse, error) {
	// Check user role in company
	role, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	if !role.Role.HasPermission(user.RoleManager) {
		return nil, errors.NewForbiddenError("you need at least manager role to view SMTP configs")
	}

	// Get all SMTP configs for company
	configs, err := s.smtpRepo.FindByCompanyID(companyID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch SMTP configs", err)
	}

	var responses []*SMTPConfigResponse
	for _, config := range configs {
		responses = append(responses, s.toResponse(config))
	}

	return &ListSMTPConfigsResponse{
		Configs: responses,
	}, nil
}

func (s *Service) GetSMTPConfig(userID, companyID, configID uint) (*SMTPConfigResponse, error) {
	// Check user role in company
	role, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	if !role.Role.HasPermission(user.RoleManager) {
		return nil, errors.NewForbiddenError("you need at least manager role to view SMTP configs")
	}

	// Get SMTP config
	config, err := s.smtpRepo.FindByIDAndCompany(configID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("SMTP config not found")
	}

	return s.toResponse(config), nil
}

func (s *Service) UpdateSMTPConfig(userID, companyID, configID uint, req *UpdateSMTPConfigRequest) (*SMTPConfigResponse, error) {
	// Check user role in company
	role, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	if role.Role != user.RoleOwner && role.Role != user.RoleAdmin {
		return nil, errors.NewForbiddenError("only owners and admins can update SMTP configs")
	}

	// Get existing config
	config, err := s.smtpRepo.FindByIDAndCompany(configID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("SMTP config not found")
	}

	// Update fields
	if req.Host != nil {
		config.Host = *req.Host
	}
	if req.User != nil {
		config.User = *req.User
	}
	if req.Password != nil {
		// Encrypt new password (using AES so it can be decrypted for SMTP)
		encryptedPassword, err := encryption.Encrypt(*req.Password)
		if err != nil {
			return nil, errors.NewInternalError("failed to encrypt password", err)
		}
		config.Password = encryptedPassword
	}
	if req.Port != nil {
		if *req.Port < 1 || *req.Port > 65535 {
			return nil, errors.NewValidationError("port must be between 1 and 65535")
		}
		config.Port = *req.Port
	}
	if req.FromName != nil {
		config.FromName = *req.FromName
	}
	if req.Security != nil {
		securityType := smtpconfig.SecurityType(*req.Security)
		if !securityType.IsValid() {
			return nil, errors.NewValidationError("invalid security type")
		}
		config.Security = securityType
	}
	if req.RateLimit != nil {
		if *req.RateLimit < 1 {
			return nil, errors.NewValidationError("rate limit must be at least 1")
		}
		config.RateLimit = *req.RateLimit
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	}

	if err := s.smtpRepo.Update(config); err != nil {
		return nil, errors.NewInternalError("failed to update SMTP config", err)
	}

	return s.toResponse(config), nil
}

func (s *Service) DeleteSMTPConfig(userID, companyID, configID uint) error {
	// Check user role in company
	role, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return errors.NewForbiddenError("you don't have access to this company")
	}

	if role.Role != user.RoleOwner && role.Role != user.RoleAdmin {
		return errors.NewForbiddenError("only owners and admins can delete SMTP configs")
	}

	// Verify config belongs to company
	_, err = s.smtpRepo.FindByIDAndCompany(configID, companyID)
	if err != nil {
		return errors.NewNotFoundError("SMTP config not found")
	}

	if err := s.smtpRepo.Delete(configID); err != nil {
		return errors.NewInternalError("failed to delete SMTP config", err)
	}

	return nil
}

func (s *Service) SetDefaultSMTPConfig(userID, companyID, configID uint) error {
	// Check user role in company
	role, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return errors.NewForbiddenError("you don't have access to this company")
	}

	if role.Role != user.RoleOwner && role.Role != user.RoleAdmin {
		return errors.NewForbiddenError("only owners and admins can set default SMTP config")
	}

	// Verify config belongs to company
	_, err = s.smtpRepo.FindByIDAndCompany(configID, companyID)
	if err != nil {
		return errors.NewNotFoundError("SMTP config not found")
	}

	if err := s.smtpRepo.SetAsDefault(configID, companyID); err != nil {
		return errors.NewInternalError("failed to set default SMTP config", err)
	}

	return nil
}

func (s *Service) toResponse(config *smtpconfig.SMTPConfig) *SMTPConfigResponse {
	return &SMTPConfigResponse{
		ID:        config.ID,
		CompanyID: config.CompanyID,
		Host:      config.Host,
		User:      config.User,
		Port:      config.Port,
		FromName:  config.FromName,
		Security:  config.Security.String(),
		RateLimit: config.RateLimit,
		IsActive:  config.IsActive,
		IsDefault: config.IsDefault,
		CreatedAt: config.CreatedAt.Format(time.RFC3339),
		UpdatedAt: config.UpdatedAt.Format(time.RFC3339),
	}
}

