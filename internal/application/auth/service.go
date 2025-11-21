package auth

import (
	"fmt"
	"os"

	emailApp "github.com/YasserCherfaoui/darween/internal/application/email"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/internal/infrastructure/security"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	userRepo      user.Repository
	jwtManager    *security.JWTManager
	emailService  *emailApp.Service
}

func NewService(userRepo user.Repository, jwtManager *security.JWTManager, emailService *emailApp.Service) *Service {
	return &Service{
		userRepo:     userRepo,
		jwtManager:   jwtManager,
		emailService: emailService,
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
	if len(userCompanies) > 0 {
		companyID = userCompanies[0].CompanyID
	}

	// Build reset URL
	resetURL := os.Getenv("FRONTEND_URL")
	if resetURL == "" {
		resetURL = "http://localhost:3000"
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
