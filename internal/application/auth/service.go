package auth

import (
	"fmt"
	"os"

	emailApp "github.com/YasserCherfaoui/darween/internal/application/email"
	otpApp "github.com/YasserCherfaoui/darween/internal/application/otp"
	"github.com/YasserCherfaoui/darween/internal/domain/company"
	invitationDomain "github.com/YasserCherfaoui/darween/internal/domain/invitation"
	otpDomain "github.com/YasserCherfaoui/darween/internal/domain/otp"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/internal/infrastructure/security"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	userRepo       user.Repository
	companyRepo    company.Repository
	invitationRepo invitationDomain.Repository
	jwtManager     *security.JWTManager
	emailService   *emailApp.Service
	otpService     *otpApp.Service
}

func NewService(userRepo user.Repository, companyRepo company.Repository, invitationRepo invitationDomain.Repository, jwtManager *security.JWTManager, emailService *emailApp.Service, otpService *otpApp.Service) *Service {
	return &Service{
		userRepo:       userRepo,
		companyRepo:    companyRepo,
		invitationRepo: invitationRepo,
		jwtManager:     jwtManager,
		emailService:   emailService,
		otpService:     otpService,
	}
}

func (s *Service) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.NewConflictError("user with this email already exists")
	}

	// Create new user
	newUser := &user.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	// Hash password
	if err := newUser.HashPassword(req.Password); err != nil {
		return nil, errors.NewInternalError("failed to hash password", err)
	}

	// Save user
	if err := s.userRepo.Create(newUser); err != nil {
		return nil, errors.NewInternalError("failed to create user", err)
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(newUser.ID, newUser.Email)
	if err != nil {
		return nil, errors.NewInternalError("failed to generate token", err)
	}

	return &AuthResponse{
		Token: token,
		User: UserInfo{
			ID:        newUser.ID,
			Email:     newUser.Email,
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
		},
	}, nil
}

func (s *Service) Login(req *LoginRequest) (*AuthResponse, error) {
	// Find user by email
	u, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid email or password")
	}

	// Check if user is active
	if !u.IsActive {
		return nil, errors.NewUnauthorizedError("user account is inactive")
	}

	// Verify password
	if !u.CheckPassword(req.Password) {
		return nil, errors.NewUnauthorizedError("invalid email or password")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(u.ID, u.Email)
	if err != nil {
		return nil, errors.NewInternalError("failed to generate token", err)
	}

	return &AuthResponse{
		Token: token,
		User: UserInfo{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		},
	}, nil
}

// RequestPasswordReset generates a password reset token and sends an email
func (s *Service) RequestPasswordReset(req *PasswordResetRequest) (*PasswordResetResponse, error) {
	// Find user by email
	u, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Don't reveal if user exists or not for security
		return &PasswordResetResponse{
			Message: "If an account with that email exists, a password reset link has been sent.",
		}, nil
	}

	// Check if user is active
	if !u.IsActive {
		// Still return success to avoid revealing account status
		return &PasswordResetResponse{
			Message: "If an account with that email exists, a password reset link has been sent.",
		}, nil
	}

	// Generate password reset token
	token, err := s.jwtManager.GeneratePasswordResetToken(u.ID, u.Email)
	if err != nil {
		return nil, errors.NewInternalError("failed to generate reset token", err)
	}

	// Get user's first company (if any) for email context
	userCompanies, _ := s.userRepo.FindUserCompaniesByUserID(u.ID)
	companyID := uint(0)
	var companyERPUrl string
	if len(userCompanies) > 0 {
		companyID = userCompanies[0].CompanyID
		// Get company to access ERPUrl
		if company, err := s.companyRepo.FindByID(companyID); err == nil {
			companyERPUrl = company.ERPUrl
		}
	}

	// Build reset URL from company ERPUrl or fallback to env/default
	resetURL := companyERPUrl
	if resetURL == "" {
		resetURL = os.Getenv("FRONTEND_URL")
	}
	if resetURL == "" {
		resetURL = "http://localhost:3000"
	}
	// Ensure URL doesn't end with /
	if len(resetURL) > 0 && resetURL[len(resetURL)-1] == '/' {
		resetURL = resetURL[:len(resetURL)-1]
	}
	resetURL = fmt.Sprintf("%s/reset-password?token=%s", resetURL, token)

	// Send password reset email
	userName := u.FirstName
	if userName == "" {
		userName = u.Email
	}

	emailReq := &emailApp.SendPasswordResetEmailRequest{
		CompanyID:  companyID,
		UserEmail:  u.Email,
		ResetToken: token,
		ResetURL:   resetURL,
		UserName:   userName,
	}

	// Send email (don't fail if email sending fails, just log it)
	if err := s.emailService.SendPasswordResetEmail(emailReq); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to send password reset email: %v\n", err)
	}

	return &PasswordResetResponse{
		Message: "If an account with that email exists, a password reset link has been sent.",
	}, nil
}

// ConfirmPasswordReset validates the token and resets the password
func (s *Service) ConfirmPasswordReset(req *PasswordResetConfirmRequest) (*PasswordResetResponse, error) {
	// Validate token
	claims, err := s.jwtManager.ValidateToken(req.Token)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid or expired reset token")
	}

	// Find user
	u, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	// Verify email matches
	if u.Email != claims.Email {
		return nil, errors.NewUnauthorizedError("invalid token")
	}

	// Check if user is active
	if !u.IsActive {
		return nil, errors.NewUnauthorizedError("user account is inactive")
	}

	// Hash new password
	if err := u.HashPassword(req.Password); err != nil {
		return nil, errors.NewInternalError("failed to hash password", err)
	}

	// Update user
	if err := s.userRepo.Update(u); err != nil {
		return nil, errors.NewInternalError("failed to update password", err)
	}

	return &PasswordResetResponse{
		Message: "Password has been reset successfully.",
	}, nil
}

// ValidateInvitation validates an invitation token and returns the email and company ID
func (s *Service) ValidateInvitation(req *ValidateInvitationRequest) (*ValidateInvitationResponse, error) {
	inv, err := s.invitationRepo.FindByToken(req.Token)
	if err != nil {
		return &ValidateInvitationResponse{
			Valid: false,
		}, nil
	}

	return &ValidateInvitationResponse{
		Email:     inv.Email,
		CompanyID: inv.CompanyID,
		Valid:     true,
	}, nil
}

// AcceptInvitation accepts an invitation and sets up the user account
func (s *Service) AcceptInvitation(req *AcceptInvitationRequest) (*AuthResponse, error) {
	// Validate token from database
	inv, err := s.invitationRepo.FindByToken(req.Token)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid or expired invitation token")
	}

	email := inv.Email

	// Find user by email
	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	// Check if user is already set up (has first name and last name)
	// Users created automatically have empty first/last name, so we allow them to complete setup
	if u.FirstName != "" && u.LastName != "" {
		// User already has name set, check if they can login (password is set)
		// If they can't login with a dummy password, they need to set password
		// If they can login, they're already set up
		if u.Password != "" {
			// Try a dummy check - if it doesn't error, password might be set
			// We'll allow password update if they provide the correct current password
			// For now, we'll allow them to update their password via invitation
		}
	}

	// Update user with first name, last name, and password
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	if err := u.HashPassword(req.Password); err != nil {
		return nil, errors.NewInternalError("failed to hash password", err)
	}

	// Update user
	if err := s.userRepo.Update(u); err != nil {
		return nil, errors.NewInternalError("failed to update user", err)
	}

	// Ensure user is active
	u.IsActive = true
	if err := s.userRepo.Update(u); err != nil {
		return nil, errors.NewInternalError("failed to activate user", err)
	}

	// Mark invitation as used
	if err := s.invitationRepo.MarkAsUsed(req.Token); err != nil {
		// Log but don't fail - user is already set up
		fmt.Printf("Failed to mark invitation as used: %v\n", err)
	}

	// Generate JWT token for login
	token, err := s.jwtManager.GenerateToken(u.ID, u.Email)
	if err != nil {
		return nil, errors.NewInternalError("failed to generate token", err)
	}

	return &AuthResponse{
		Token: token,
		User: UserInfo{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		},
	}, nil
}

// ValidateOTP validates an OTP code for a given email
func (s *Service) ValidateOTP(req *ValidateOTPRequest) (*ValidateOTPResponse, error) {
	codeStr := req.GetCodeString()
	otpRecord, err := s.otpService.ValidateOTP(codeStr, req.Email)
	if err != nil {
		return &ValidateOTPResponse{
			Valid: false,
		}, nil
	}

	return &ValidateOTPResponse{
		Valid:   true,
		Purpose: string(otpRecord.Purpose),
		UserID:  otpRecord.UserID,
	}, nil
}

// CompleteUserSetup completes user setup with OTP (for new users)
func (s *Service) CompleteUserSetup(req *CompleteUserSetupRequest) (*AuthResponse, error) {
	// Validate OTP
	codeStr := req.GetCodeString()
	otpRecord, err := s.otpService.ValidateOTP(codeStr, req.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid or expired OTP code")
	}

	// Check OTP purpose
	if otpRecord.Purpose != otpDomain.PurposeSetupAccount {
		return nil, errors.NewValidationError("OTP is not for account setup")
	}

	// Find user
	u, err := s.userRepo.FindByID(otpRecord.UserID)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	// Verify email matches
	if u.Email != req.Email {
		return nil, errors.NewUnauthorizedError("email does not match OTP")
	}

	// Update user with provided information
	u.FirstName = req.FirstName
	u.LastName = req.LastName

	// Hash new password
	if err := u.HashPassword(req.Password); err != nil {
		return nil, errors.NewInternalError("failed to hash password", err)
	}

	// Ensure user is active
	u.IsActive = true

	// Update user
	if err := s.userRepo.Update(u); err != nil {
		return nil, errors.NewInternalError("failed to update user", err)
	}

	// Mark OTP as used
	if err := s.otpService.MarkAsUsed(otpRecord.ID); err != nil {
		// Log but don't fail - user is already set up
		fmt.Printf("Failed to mark OTP as used: %v\n", err)
	}

	// Generate JWT token for login
	token, err := s.jwtManager.GenerateToken(u.ID, u.Email)
	if err != nil {
		return nil, errors.NewInternalError("failed to generate token", err)
	}

	return &AuthResponse{
		Token: token,
		User: UserInfo{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		},
	}, nil
}

// ChangePasswordWithOTP changes user password using OTP
func (s *Service) ChangePasswordWithOTP(req *ChangePasswordWithOTPRequest) (*PasswordResetResponse, error) {
	// Validate OTP
	codeStr := req.GetCodeString()
	otpRecord, err := s.otpService.ValidateOTP(codeStr, req.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid or expired OTP code")
	}

	// Check OTP purpose
	if otpRecord.Purpose != otpDomain.PurposePasswordReset {
		return nil, errors.NewValidationError("OTP is not for password reset")
	}

	// Find user
	u, err := s.userRepo.FindByID(otpRecord.UserID)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	// Verify email matches
	if u.Email != req.Email {
		return nil, errors.NewUnauthorizedError("email does not match OTP")
	}

	// Check if user is active
	if !u.IsActive {
		return nil, errors.NewUnauthorizedError("user account is inactive")
	}

	// Hash new password
	if err := u.HashPassword(req.Password); err != nil {
		return nil, errors.NewInternalError("failed to hash password", err)
	}

	// Update user
	if err := s.userRepo.Update(u); err != nil {
		return nil, errors.NewInternalError("failed to update password", err)
	}

	// Mark OTP as used
	if err := s.otpService.MarkAsUsed(otpRecord.ID); err != nil {
		// Log but don't fail - password is already updated
		fmt.Printf("Failed to mark OTP as used: %v\n", err)
	}

	return &PasswordResetResponse{
		Message: "Password has been changed successfully.",
	}, nil
}
