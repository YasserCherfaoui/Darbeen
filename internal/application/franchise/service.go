package franchise

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	emailApp "github.com/YasserCherfaoui/darween/internal/application/email"
	otpApp "github.com/YasserCherfaoui/darween/internal/application/otp"
	companyDomain "github.com/YasserCherfaoui/darween/internal/domain/company"
	franchiseDomain "github.com/YasserCherfaoui/darween/internal/domain/franchise"
	inventoryDomain "github.com/YasserCherfaoui/darween/internal/domain/inventory"
	invitationDomain "github.com/YasserCherfaoui/darween/internal/domain/invitation"
	otpDomain "github.com/YasserCherfaoui/darween/internal/domain/otp"
	productDomain "github.com/YasserCherfaoui/darween/internal/domain/product"
	userDomain "github.com/YasserCherfaoui/darween/internal/domain/user"
	smtpconfigDomain "github.com/YasserCherfaoui/darween/internal/domain/smtpconfig"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	franchiseRepo franchiseDomain.Repository
	inventoryRepo inventoryDomain.Repository
	companyRepo   companyDomain.Repository
	userRepo      userDomain.Repository
	productRepo   productDomain.Repository
	emailService  *emailApp.Service
	smtpRepo      smtpconfigDomain.Repository
	invitationRepo invitationDomain.Repository
	otpService    *otpApp.Service
}

func NewService(
	franchiseRepo franchiseDomain.Repository,
	inventoryRepo inventoryDomain.Repository,
	companyRepo companyDomain.Repository,
	userRepo userDomain.Repository,
	productRepo productDomain.Repository,
	emailService *emailApp.Service,
	smtpRepo smtpconfigDomain.Repository,
	invitationRepo invitationDomain.Repository,
	otpService *otpApp.Service,
) *Service {
	return &Service{
		franchiseRepo: franchiseRepo,
		inventoryRepo: inventoryRepo,
		companyRepo:   companyRepo,
		userRepo:      userRepo,
		productRepo:   productRepo,
		emailService:  emailService,
		smtpRepo:      smtpRepo,
		invitationRepo: invitationRepo,
		otpService:    otpService,
	}
}

func (s *Service) CreateFranchise(userID, companyID uint, req *CreateFranchiseRequest) (*FranchiseResponse, error) {
	// Check user is owner/admin of the parent company
	role, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	if role.Role != userDomain.RoleOwner && role.Role != userDomain.RoleAdmin {
		return nil, errors.NewForbiddenError("only owners and admins can create franchises")
	}

	// Check if franchise code already exists
	existingFranchise, _ := s.franchiseRepo.FindByCode(req.Code)
	if existingFranchise != nil {
		return nil, errors.NewConflictError("franchise with this code already exists")
	}

	// Create franchise
	newFranchise := &franchiseDomain.Franchise{
		ParentCompanyID: companyID,
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		Address:         req.Address,
		Phone:           req.Phone,
		Email:           req.Email,
		IsActive:        true,
	}

	if err := s.franchiseRepo.Create(newFranchise); err != nil {
		return nil, errors.NewInternalError("failed to create franchise", err)
	}

	return s.buildFranchiseResponse(newFranchise), nil
}

func (s *Service) GetFranchisesByCompanyID(userID, companyID uint) ([]*FranchiseResponse, error) {
	// Check user has access to company
	_, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	franchises, err := s.franchiseRepo.FindByParentCompanyID(companyID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch franchises", err)
	}

	result := make([]*FranchiseResponse, 0, len(franchises))
	for _, f := range franchises {
		result = append(result, s.buildFranchiseResponse(f))
	}

	return result, nil
}

func (s *Service) GetFranchiseByID(userID, franchiseID uint) (*FranchiseResponse, error) {
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	// Check if user is parent company admin or franchise user
	parentRole, _ := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
	franchiseRole, _ := s.userRepo.FindUserRoleInFranchise(userID, franchiseID)

	if (parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin)) &&
		franchiseRole == nil {
		return nil, errors.NewForbiddenError("you don't have access to this franchise")
	}

	return s.buildFranchiseResponse(franchise), nil
}

func (s *Service) UpdateFranchise(userID, franchiseID uint, req *UpdateFranchiseRequest) (*FranchiseResponse, error) {
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	// Check if user is parent company admin or franchise admin
	parentRole, _ := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
	franchiseRole, _ := s.userRepo.FindUserRoleInFranchise(userID, franchiseID)

	if (parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin)) &&
		(franchiseRole == nil || (franchiseRole.Role != userDomain.RoleOwner && franchiseRole.Role != userDomain.RoleAdmin)) {
		return nil, errors.NewForbiddenError("only owners and admins can update franchise")
	}

	// Update fields
	if req.Name != "" {
		franchise.Name = req.Name
	}
	if req.Description != "" {
		franchise.Description = req.Description
	}
	if req.Address != "" {
		franchise.Address = req.Address
	}
	if req.Phone != "" {
		franchise.Phone = req.Phone
	}
	if req.Email != "" {
		franchise.Email = req.Email
	}
	if req.IsActive != nil {
		franchise.IsActive = *req.IsActive
	}

	if err := s.franchiseRepo.Update(franchise); err != nil {
		return nil, errors.NewInternalError("failed to update franchise", err)
	}

	return s.buildFranchiseResponse(franchise), nil
}

func (s *Service) InitializeFranchiseInventory(userID, franchiseID uint) error {
	// Verify franchise exists and user has access
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return errors.NewNotFoundError("franchise not found")
	}

	parentRole, _ := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
	if parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin) {
		return errors.NewForbiddenError("only parent company owners and admins can initialize franchise inventory")
	}

	// Get all products from parent company
	products, _, err := s.productRepo.FindProductsByCompanyID(franchise.ParentCompanyID, 1, 1000)
	if err != nil {
		return errors.NewInternalError("failed to fetch products", err)
	}

	// Create inventory records for all variants
	for _, product := range products {
		variants, err := s.productRepo.FindProductVariantsByProductID(product.ID)
		if err != nil {
			continue
		}

		for _, variant := range variants {
			// Check if inventory already exists
			existing, _ := s.inventoryRepo.FindByVariantAndFranchise(variant.ID, franchiseID)
			if existing == nil {
				newInventory := &inventoryDomain.Inventory{
					ProductVariantID: variant.ID,
					CompanyID:        nil,
					FranchiseID:      uintPtr(franchiseID),
					Stock:            0,
					ReservedStock:    0,
					IsActive:         true,
				}

				if err := s.inventoryRepo.Create(newInventory); err != nil {
					// Log error but continue
					continue
				}

				// Log initial movement
				movement := &inventoryDomain.InventoryMovement{
					InventoryID:   newInventory.ID,
					MovementType:  inventoryDomain.MovementTypeAdjustment,
					Quantity:      0,
					PreviousStock: 0,
					NewStock:      0,
					Notes:         stringPtr("Initialized from catalog"),
					CreatedByID:   userID,
				}
				s.inventoryRepo.CreateMovement(movement)
			}
		}
	}

	return nil
}

func (s *Service) SetFranchisePricing(userID, franchiseID uint, req *SetFranchisePricingRequest) (*FranchisePricingResponse, error) {
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	// Only parent company owner/admin can set pricing
	parentRole, _ := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
	if parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin) {
		return nil, errors.NewForbiddenError("only parent company owners and admins can set franchise pricing")
	}

	// Verify variant exists and belongs to parent company
	variant, err := s.productRepo.FindProductVariantByID(req.ProductVariantID)
	if err != nil {
		return nil, errors.NewNotFoundError("product variant not found")
	}

	product, err := s.productRepo.FindProductByID(variant.ProductID)
	if err != nil || product.CompanyID != franchise.ParentCompanyID {
		return nil, errors.NewForbiddenError("product variant does not belong to parent company")
	}

	// Check if pricing already exists
	existingPricing, _ := s.franchiseRepo.FindPricing(franchiseID, req.ProductVariantID)

	var pricing *franchiseDomain.FranchisePricing
	if existingPricing != nil {
		pricing = existingPricing
		if req.RetailPrice != nil {
			pricing.SetRetailPriceOverride(*req.RetailPrice)
		}
		if req.WholesalePrice != nil {
			pricing.SetWholesalePriceOverride(*req.WholesalePrice)
		}
		if err := s.franchiseRepo.UpdatePricing(pricing); err != nil {
			return nil, errors.NewInternalError("failed to update pricing", err)
		}
	} else {
		pricing = &franchiseDomain.FranchisePricing{
			FranchiseID:      franchiseID,
			ProductVariantID: req.ProductVariantID,
			RetailPrice:      req.RetailPrice,
			WholesalePrice:   req.WholesalePrice,
			IsActive:         true,
		}
		if err := s.franchiseRepo.CreatePricing(pricing); err != nil {
			return nil, errors.NewInternalError("failed to create pricing", err)
		}
	}

	return s.buildFranchisePricingResponse(pricing, product, variant), nil
}

func (s *Service) GetFranchisePricing(userID, franchiseID uint) ([]*FranchisePricingResponse, error) {
	// Verify access
	_, err := s.GetFranchiseByID(userID, franchiseID)
	if err != nil {
		return nil, err
	}

	pricings, err := s.franchiseRepo.FindAllPricingByFranchise(franchiseID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch pricing", err)
	}

	result := make([]*FranchisePricingResponse, 0, len(pricings))
	for _, pricing := range pricings {
		variant, _ := s.productRepo.FindProductVariantByID(pricing.ProductVariantID)
		if variant == nil {
			continue
		}
		product, _ := s.productRepo.FindProductByID(variant.ProductID)
		if product == nil {
			continue
		}

		result = append(result, s.buildFranchisePricingResponse(pricing, product, variant))
	}

	return result, nil
}

func (s *Service) DeleteFranchisePricing(userID, franchiseID, variantID uint) error {
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return errors.NewNotFoundError("franchise not found")
	}

	// Only parent company owner/admin can delete pricing
	parentRole, _ := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
	if parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin) {
		return errors.NewForbiddenError("only parent company owners and admins can delete franchise pricing")
	}

	if err := s.franchiseRepo.DeletePricing(franchiseID, variantID); err != nil {
		return errors.NewInternalError("failed to delete pricing", err)
	}

	return nil
}

func (s *Service) BulkSetFranchisePricing(userID, franchiseID uint, req *BulkSetFranchisePricingRequest) (*BulkSetFranchisePricingResponse, error) {
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	// Only parent company owner/admin can set pricing
	parentRole, _ := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
	if parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin) {
		return nil, errors.NewForbiddenError("only parent company owners and admins can set franchise pricing")
	}

	// Verify product exists and belongs to parent company
	product, err := s.productRepo.FindProductByID(req.ProductID)
	if err != nil {
		return nil, errors.NewNotFoundError("product not found")
	}

	if product.CompanyID != franchise.ParentCompanyID {
		return nil, errors.NewForbiddenError("product does not belong to parent company")
	}

	// Get all variants for the product
	variants, err := s.productRepo.FindProductVariantsByProductID(req.ProductID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch product variants", err)
	}

	result := make([]*FranchisePricingResponse, 0, len(variants))
	updatedCount := 0

	// Set pricing for each variant
	for _, variant := range variants {
		// Check if pricing already exists
		existingPricing, _ := s.franchiseRepo.FindPricing(franchiseID, variant.ID)

		var pricing *franchiseDomain.FranchisePricing
		if existingPricing != nil {
			pricing = existingPricing
			if req.RetailPrice != nil {
				pricing.SetRetailPriceOverride(*req.RetailPrice)
			}
			if req.WholesalePrice != nil {
				pricing.SetWholesalePriceOverride(*req.WholesalePrice)
			}
			if err := s.franchiseRepo.UpdatePricing(pricing); err != nil {
				continue // Skip this variant on error
			}
		} else {
			pricing = &franchiseDomain.FranchisePricing{
				FranchiseID:      franchiseID,
				ProductVariantID: variant.ID,
				RetailPrice:      req.RetailPrice,
				WholesalePrice:   req.WholesalePrice,
				IsActive:         true,
			}
			if err := s.franchiseRepo.CreatePricing(pricing); err != nil {
				continue // Skip this variant on error
			}
		}

		result = append(result, s.buildFranchisePricingResponse(pricing, product, variant))
		updatedCount++
	}

	return &BulkSetFranchisePricingResponse{
		UpdatedCount: updatedCount,
		Pricing:      result,
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

func (s *Service) AddUserToFranchise(requestUserID, franchiseID uint, req *AddUserToFranchiseRequest) (*AddUserToFranchiseResponse, error) {
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	// Check if requester has permission (parent company admin or franchise admin)
	parentRole, _ := s.userRepo.FindUserRoleInCompany(requestUserID, franchise.ParentCompanyID)
	franchiseRole, _ := s.userRepo.FindUserRoleInFranchise(requestUserID, franchiseID)

	if (parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin)) &&
		(franchiseRole == nil || (franchiseRole.Role != userDomain.RoleOwner && franchiseRole.Role != userDomain.RoleAdmin)) {
		return nil, errors.NewForbiddenError("only owners and admins can add users to franchise")
	}

	// Validate role
	role := userDomain.Role(req.Role)
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

		targetUser = &userDomain.User{
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

	// Check if user already belongs to franchise
	existingRole, _ := s.userRepo.FindUserRoleInFranchise(targetUser.ID, franchiseID)
	if existingRole != nil {
		return nil, errors.NewConflictError("user already belongs to this franchise")
	}

	// Create user-franchise relationship
	ufr := &userDomain.UserFranchiseRole{
		UserID:      targetUser.ID,
		FranchiseID: franchiseID,
		Role:        role,
		IsActive:    true,
	}

	if err := s.userRepo.CreateUserFranchiseRole(ufr); err != nil {
		return nil, errors.NewInternalError("failed to add user to franchise", err)
	}

	// Get company info for email
	company, err := s.companyRepo.FindByID(franchise.ParentCompanyID)
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
	defaultSMTPConfig, _ := s.smtpRepo.FindDefaultByCompanyID(franchise.ParentCompanyID)
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
			franchiseIDPtr := &franchiseID
			otpCode, err := s.otpService.GenerateOTP(
				targetUser.Email,
				targetUser.ID,
				franchise.ParentCompanyID,
				franchiseIDPtr,
				otpDomain.PurposeSetupAccount,
				24*time.Hour, // 24 hours expiry
			)
			if err == nil {
				setupURL := fmt.Sprintf("%s/setup-account?otp=%s&email=%s", baseURL, otpCode, targetUser.Email)

				// Send new user setup email
				emailReq := &emailApp.SendNewUserSetupEmailRequest{
					CompanyID:   franchise.ParentCompanyID,
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
				CompanyID:   franchise.ParentCompanyID,
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
	response := &AddUserToFranchiseResponse{
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

func (s *Service) RemoveUserFromFranchise(requestUserID, franchiseID, targetUserID uint) error {
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return errors.NewNotFoundError("franchise not found")
	}

	// Check if requester has permission
	parentRole, _ := s.userRepo.FindUserRoleInCompany(requestUserID, franchise.ParentCompanyID)
	franchiseRole, _ := s.userRepo.FindUserRoleInFranchise(requestUserID, franchiseID)

	if (parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin)) &&
		(franchiseRole == nil || (franchiseRole.Role != userDomain.RoleOwner && franchiseRole.Role != userDomain.RoleAdmin)) {
		return errors.NewForbiddenError("only owners and admins can remove users from franchise")
	}

	// Check if target user is an owner (owners cannot be removed)
	targetRole, err := s.userRepo.FindUserRoleInFranchise(targetUserID, franchiseID)
	if err != nil {
		return errors.NewNotFoundError("user not found in franchise")
	}

	if targetRole.Role == userDomain.RoleOwner {
		return errors.NewForbiddenError("cannot remove franchise owner")
	}

	if err := s.userRepo.DeleteUserFranchiseRole(targetUserID, franchiseID); err != nil {
		return errors.NewInternalError("failed to remove user from franchise", err)
	}

	return nil
}

func (s *Service) GetFranchiseUsers(requestUserID, franchiseID uint) (*ListFranchiseUsersResponse, error) {
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	// Check if requester has access (parent company or franchise role)
	parentRole, _ := s.userRepo.FindUserRoleInCompany(requestUserID, franchise.ParentCompanyID)
	franchiseRole, _ := s.userRepo.FindUserRoleInFranchise(requestUserID, franchiseID)

	if parentRole == nil && franchiseRole == nil {
		return nil, errors.NewForbiddenError("you don't have access to this franchise")
	}

	// Check if requester has at least manager role
	hasManagerAccess := false
	if parentRole != nil && parentRole.Role.HasPermission(userDomain.RoleManager) {
		hasManagerAccess = true
	}
	if franchiseRole != nil && franchiseRole.Role.HasPermission(userDomain.RoleManager) {
		hasManagerAccess = true
	}

	if !hasManagerAccess {
		return nil, errors.NewForbiddenError("you need at least manager role to view franchise users")
	}

	// Get all users in the franchise
	roles, err := s.userRepo.FindFranchiseUsersByFranchiseID(franchiseID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch franchise users", err)
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

	return &ListFranchiseUsersResponse{
		Users: users,
	}, nil
}

func (s *Service) UpdateUserRoleInFranchise(requestUserID, franchiseID, targetUserID uint, req *UpdateUserRoleRequest) error {
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return errors.NewNotFoundError("franchise not found")
	}

	// Check if requester has permission (parent company admin or franchise admin)
	parentRole, _ := s.userRepo.FindUserRoleInCompany(requestUserID, franchise.ParentCompanyID)
	franchiseRole, _ := s.userRepo.FindUserRoleInFranchise(requestUserID, franchiseID)

	if (parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin)) &&
		(franchiseRole == nil || (franchiseRole.Role != userDomain.RoleOwner && franchiseRole.Role != userDomain.RoleAdmin)) {
		return errors.NewForbiddenError("only owners and admins can update user roles")
	}

	// Validate new role
	newRole := userDomain.Role(req.Role)
	if !newRole.IsValid() {
		return errors.NewValidationError("invalid role")
	}

	// Check if target user exists in franchise
	targetRole, err := s.userRepo.FindUserRoleInFranchise(targetUserID, franchiseID)
	if err != nil {
		return errors.NewNotFoundError("user not found in franchise")
	}

	// Cannot change owner role
	if targetRole.Role == userDomain.RoleOwner {
		return errors.NewForbiddenError("cannot change owner role")
	}

	// Cannot assign owner role (only through franchise creation)
	if newRole == userDomain.RoleOwner {
		return errors.NewForbiddenError("cannot assign owner role")
	}

	// Update role
	targetRole.Role = newRole
	if err := s.userRepo.UpdateUserFranchiseRole(targetRole); err != nil {
		return errors.NewInternalError("failed to update user role", err)
	}

	return nil
}

// Helper methods

func (s *Service) buildFranchiseResponse(f *franchiseDomain.Franchise) *FranchiseResponse {
	return &FranchiseResponse{
		ID:              f.ID,
		ParentCompanyID: f.ParentCompanyID,
		Name:            f.Name,
		Code:            f.Code,
		Description:     f.Description,
		Address:         f.Address,
		Phone:           f.Phone,
		Email:           f.Email,
		IsActive:        f.IsActive,
		CreatedAt:       f.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       f.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *Service) buildFranchisePricingResponse(pricing *franchiseDomain.FranchisePricing, product *productDomain.Product, variant *productDomain.ProductVariant) *FranchisePricingResponse {
	// Get default prices from variant (with fallback to product)
	var defaultRetailPrice, defaultWholesalePrice float64

	if variant.RetailPrice != nil {
		defaultRetailPrice = *variant.RetailPrice
	} else if product.BaseRetailPrice > 0 {
		defaultRetailPrice = product.BaseRetailPrice
	}

	if variant.WholesalePrice != nil {
		defaultWholesalePrice = *variant.WholesalePrice
	} else if product.BaseWholesalePrice > 0 {
		defaultWholesalePrice = product.BaseWholesalePrice
	}

	return &FranchisePricingResponse{
		ID:                    pricing.ID,
		FranchiseID:           pricing.FranchiseID,
		ProductVariantID:      pricing.ProductVariantID,
		VariantName:           variant.Name,
		VariantSKU:            variant.SKU,
		RetailPrice:           pricing.RetailPrice,
		WholesalePrice:        pricing.WholesalePrice,
		DefaultRetailPrice:    defaultRetailPrice,
		DefaultWholesalePrice: defaultWholesalePrice,
		IsActive:              pricing.IsActive,
	}
}

// Helper functions

func uintPtr(u uint) *uint {
	return &u
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
