package company

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	emailApp "github.com/YasserCherfaoui/darween/internal/application/email"
	otpApp "github.com/YasserCherfaoui/darween/internal/application/otp"
	"github.com/YasserCherfaoui/darween/internal/domain/company"
	invitationDomain "github.com/YasserCherfaoui/darween/internal/domain/invitation"
	otpDomain "github.com/YasserCherfaoui/darween/internal/domain/otp"
	"github.com/YasserCherfaoui/darween/internal/domain/subscription"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	smtpconfigDomain "github.com/YasserCherfaoui/darween/internal/domain/smtpconfig"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	companyRepo      company.Repository
	userRepo         user.Repository
	subscriptionRepo subscription.Repository
	emailService     *emailApp.Service
	invitationRepo   invitationDomain.Repository
	smtpRepo         smtpconfigDomain.Repository
	otpService       *otpApp.Service
}

func NewService(companyRepo company.Repository, userRepo user.Repository, subscriptionRepo subscription.Repository, emailService *emailApp.Service, invitationRepo invitationDomain.Repository, smtpRepo smtpconfigDomain.Repository, otpService *otpApp.Service) *Service {
	return &Service{
		companyRepo:      companyRepo,
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
		emailService:     emailService,
		invitationRepo:   invitationRepo,
		smtpRepo:         smtpRepo,
		otpService:       otpService,
	}
}

func (s *Service) CreateCompany(userID uint, req *CreateCompanyRequest) (*CompanyResponse, error) {
	// Check if company code already exists
	existingCompany, _ := s.companyRepo.FindByCode(req.Code)
	if existingCompany != nil {
		return nil, errors.NewConflictError("company with this code already exists")
	}

	// Create company
	newCompany := &company.Company{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		ERPUrl:      req.ERPUrl,
		IsActive:    true,
	}

	if err := s.companyRepo.Create(newCompany); err != nil {
		return nil, errors.NewInternalError("failed to create company", err)
	}

	// Create user-company relationship with owner role
	ucr := &user.UserCompanyRole{
		UserID:    userID,
		CompanyID: newCompany.ID,
		Role:      user.RoleOwner,
		IsActive:  true,
	}

	if err := s.userRepo.CreateUserCompanyRole(ucr); err != nil {
		return nil, errors.NewInternalError("failed to assign owner role", err)
	}

	// Create default free subscription
	sub := &subscription.Subscription{
		CompanyID: newCompany.ID,
		PlanType:  subscription.PlanFree,
		Status:    subscription.StatusActive,
		StartDate: time.Now(),
		MaxUsers:  subscription.PlanFree.GetMaxUsers(),
	}

	if err := s.subscriptionRepo.Create(sub); err != nil {
		return nil, errors.NewInternalError("failed to create subscription", err)
	}

	return &CompanyResponse{
		ID:          newCompany.ID,
		Name:        newCompany.Name,
		Code:        newCompany.Code,
		Description: newCompany.Description,
		ERPUrl:      newCompany.ERPUrl,
		IsActive:    newCompany.IsActive,
	}, nil
}

func (s *Service) GetCompaniesByUserID(userID uint) ([]*CompanyResponse, error) {
	companies, err := s.companyRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch companies", err)
	}

	var result []*CompanyResponse
	for _, c := range companies {
		result = append(result, &CompanyResponse{
			ID:          c.ID,
			Name:        c.Name,
			Code:        c.Code,
			Description: c.Description,
			ERPUrl:      c.ERPUrl,
			IsActive:    c.IsActive,
		})
	}

	return result, nil
}

func (s *Service) GetCompanyByID(companyID uint) (*CompanyResponse, error) {
	c, err := s.companyRepo.FindByID(companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("company not found")
	}

	return &CompanyResponse{
		ID:          c.ID,
		Name:        c.Name,
		Code:        c.Code,
		Description: c.Description,
		IsActive:    c.IsActive,
	}, nil
}

func (s *Service) UpdateCompany(userID, companyID uint, req *UpdateCompanyRequest) (*CompanyResponse, error) {
	// Check user role in company
	role, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	if role.Role != user.RoleOwner && role.Role != user.RoleAdmin {
		return nil, errors.NewForbiddenError("only owners and admins can update company details")
	}

	// Get company
	c, err := s.companyRepo.FindByID(companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("company not found")
	}

	// Update fields
	if req.Name != "" {
		c.Name = req.Name
	}
	if req.Description != "" {
		c.Description = req.Description
	}
	// ERPUrl: Always update if provided (allows clearing by setting to empty string)
	// The frontend will always send this field, so we update it
	c.ERPUrl = req.ERPUrl
	if req.IsActive != nil {
		c.IsActive = *req.IsActive
	}

	if err := s.companyRepo.Update(c); err != nil {
		return nil, errors.NewInternalError("failed to update company", err)
	}

	return &CompanyResponse{
		ID:          c.ID,
		Name:        c.Name,
		Code:        c.Code,
		Description: c.Description,
		IsActive:    c.IsActive,
	}, nil
}

// generateRandomPassword generates a secure random password
func generateRandomPassword(length int) (string, error) {
	if length < 12 {
		length = 12
	}
	
	// Character sets for password generation
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		special   = "!@#$%^&*"
	)
	allChars := lowercase + uppercase + digits + special
	
	// Ensure at least one character from each set
	password := make([]byte, length)
	
	// First 4 characters: one from each set
	password[0] = lowercase[getRandomInt(len(lowercase))]
	password[1] = uppercase[getRandomInt(len(uppercase))]
	password[2] = digits[getRandomInt(len(digits))]
	password[3] = special[getRandomInt(len(special))]
	
	// Fill the rest with random characters from all sets
	for i := 4; i < length; i++ {
		password[i] = allChars[getRandomInt(len(allChars))]
	}
	
	// Shuffle the password
	for i := len(password) - 1; i > 0; i-- {
		j := getRandomInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}
	
	return string(password), nil
}

// getRandomInt returns a random integer in [0, max)
func getRandomInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func (s *Service) AddUserToCompany(requestUserID, companyID uint, req *AddUserToCompanyRequest) (*AddUserToCompanyResponse, error) {
	// Check if requester has permission
	requesterRole, err := s.userRepo.FindUserRoleInCompany(requestUserID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	if requesterRole.Role != user.RoleOwner && requesterRole.Role != user.RoleAdmin {
		return nil, errors.NewForbiddenError("only owners and admins can add users")
	}

	// Validate role
	role := user.Role(req.Role)
	if !role.IsValid() {
		return nil, errors.NewValidationError("invalid role")
	}

	// Find user by email or create if not found
	targetUser, err := s.userRepo.FindByEmail(req.Email)
	userCreated := false
	var generatedPassword string

	if err != nil {
		// User not found, create new user with random password
		generatedPassword, err = generateRandomPassword(16)
		if err != nil {
			return nil, errors.NewInternalError("failed to generate password", err)
		}

		targetUser = &user.User{
			Email:     req.Email,
			FirstName: "",
			LastName:  "",
			IsActive:  true,
		}

		if err := targetUser.HashPassword(generatedPassword); err != nil {
			return nil, errors.NewInternalError("failed to hash password", err)
		}

		if err := s.userRepo.Create(targetUser); err != nil {
			return nil, errors.NewInternalError("failed to create user", err)
		}

		userCreated = true
	}

	// Check if user already belongs to company
	existingRole, _ := s.userRepo.FindUserRoleInCompany(targetUser.ID, companyID)
	if existingRole != nil {
		return nil, errors.NewConflictError("user already belongs to this company")
	}

	// Create user-company relationship
	ucr := &user.UserCompanyRole{
		UserID:    targetUser.ID,
		CompanyID: companyID,
		Role:      role,
		IsActive:  true,
	}

	if err := s.userRepo.CreateUserCompanyRole(ucr); err != nil {
		return nil, errors.NewInternalError("failed to add user to company", err)
	}

	// Get company info
	company, err := s.companyRepo.FindByID(companyID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch company", err)
	}

	// Get inviter name
	inviter, _ := s.userRepo.FindByID(requestUserID)
	inviterName := "A team member"
	if inviter != nil {
		if inviter.FirstName != "" {
			inviterName = fmt.Sprintf("%s %s", inviter.FirstName, inviter.LastName)
		} else {
			inviterName = inviter.Email
		}
	}

	// Check if company has default SMTP config
	defaultSMTPConfig, _ := s.smtpRepo.FindDefaultByCompanyID(companyID)
	emailSent := false

	// Build base URL from company ERPUrl or fallback to env/default
	baseURL := company.ERPUrl
	if baseURL == "" {
		baseURL = os.Getenv("FRONTEND_URL")
	}
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	if defaultSMTPConfig != nil {
		if userCreated {
			// New user - generate OTP and send setup email with credentials
			otpCode, err := s.otpService.GenerateOTP(
				targetUser.Email,
				targetUser.ID,
				companyID,
				nil, // No franchise ID for company
				otpDomain.PurposeSetupAccount,
				24*time.Hour, // 24 hours expiry
			)
			if err == nil {
				setupURL := fmt.Sprintf("%s/setup-account?otp=%s&email=%s", baseURL, otpCode, targetUser.Email)

				// Send new user setup email
				emailReq := &emailApp.SendNewUserSetupEmailRequest{
					CompanyID:   companyID,
					UserEmail:   targetUser.Email,
					Password:    generatedPassword,
					InviterName: inviterName,
					OTPCode:     otpCode,
					SetupURL:    setupURL,
					CompanyName: company.Name,
				}

				if err := s.emailService.SendNewUserSetupEmail(emailReq); err != nil {
					fmt.Printf("Failed to send new user setup email: %v\n", err)
				} else {
					emailSent = true
				}
			}
		} else {
			// Existing user - send welcome email with role
			loginURL := fmt.Sprintf("%s/login", baseURL)

			emailReq := &emailApp.SendWelcomeEmailRequest{
				CompanyID:   companyID,
				UserEmail:   targetUser.Email,
				InviterName: inviterName,
				Role:        role.String(),
				LoginURL:    loginURL,
				CompanyName: company.Name,
			}

			if err := s.emailService.SendWelcomeEmail(emailReq); err != nil {
				fmt.Printf("Failed to send welcome email: %v\n", err)
			} else {
				emailSent = true
			}
		}
	}

	// Return response with credentials if user was created
	response := &AddUserToCompanyResponse{
		UserCreated: userCreated,
		EmailSent:   emailSent,
	}

	if userCreated && generatedPassword != "" {
		response.Credentials = &UserCredentials{
			Email:    targetUser.Email,
			Password: generatedPassword,
		}
	}

	return response, nil
}

func (s *Service) RemoveUserFromCompany(requestUserID, companyID, targetUserID uint) error {
	// Check if requester has permission
	requesterRole, err := s.userRepo.FindUserRoleInCompany(requestUserID, companyID)
	if err != nil {
		return errors.NewForbiddenError("you don't have access to this company")
	}

	if requesterRole.Role != user.RoleOwner && requesterRole.Role != user.RoleAdmin {
		return errors.NewForbiddenError("only owners and admins can remove users")
	}

	// Check if target user is an owner (owners cannot be removed)
	targetRole, err := s.userRepo.FindUserRoleInCompany(targetUserID, companyID)
	if err != nil {
		return errors.NewNotFoundError("user not found in company")
	}

	if targetRole.Role == user.RoleOwner {
		return errors.NewForbiddenError("cannot remove company owner")
	}

	if err := s.userRepo.DeleteUserCompanyRole(targetUserID, companyID); err != nil {
		return errors.NewInternalError("failed to remove user from company", err)
	}

	return nil
}

func (s *Service) GetCompanyUsers(requestUserID, companyID uint) (*ListCompanyUsersResponse, error) {
	// Check if requester has at least manager role
	requesterRole, err := s.userRepo.FindUserRoleInCompany(requestUserID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	if !requesterRole.Role.HasPermission(user.RoleManager) {
		return nil, errors.NewForbiddenError("you need at least manager role to view company users")
	}

	// Get all users in the company
	roles, err := s.userRepo.FindCompanyUsersByCompanyID(companyID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch company users", err)
	}

	var users []*UserWithRoleResponse
	for _, role := range roles {
		u, err := s.userRepo.FindByID(role.UserID)
		if err != nil {
			continue // Skip if user not found
		}

		users = append(users, &UserWithRoleResponse{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      role.Role.String(),
			IsActive:  role.IsActive,
		})
	}

	return &ListCompanyUsersResponse{
		Users: users,
	}, nil
}

func (s *Service) UpdateUserRoleInCompany(requestUserID, companyID, targetUserID uint, req *UpdateUserRoleRequest) error {
	// Check if requester has permission
	requesterRole, err := s.userRepo.FindUserRoleInCompany(requestUserID, companyID)
	if err != nil {
		return errors.NewForbiddenError("you don't have access to this company")
	}

	if requesterRole.Role != user.RoleOwner && requesterRole.Role != user.RoleAdmin {
		return errors.NewForbiddenError("only owners and admins can update user roles")
	}

	// Validate new role
	newRole := user.Role(req.Role)
	if !newRole.IsValid() {
		return errors.NewValidationError("invalid role")
	}

	// Check if target user exists in company
	targetRole, err := s.userRepo.FindUserRoleInCompany(targetUserID, companyID)
	if err != nil {
		return errors.NewNotFoundError("user not found in company")
	}

	// Cannot change owner role
	if targetRole.Role == user.RoleOwner {
		return errors.NewForbiddenError("cannot change owner role")
	}

	// Cannot assign owner role (only through company creation)
	if newRole == user.RoleOwner {
		return errors.NewForbiddenError("cannot assign owner role")
	}

	// Update role
	targetRole.Role = newRole
	if err := s.userRepo.UpdateUserCompanyRole(targetRole); err != nil {
		return errors.NewInternalError("failed to update user role", err)
	}

	return nil
}

