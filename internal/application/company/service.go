package company

import (
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/company"
	"github.com/YasserCherfaoui/darween/internal/domain/subscription"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	companyRepo      company.Repository
	userRepo         user.Repository
	subscriptionRepo subscription.Repository
}

func NewService(companyRepo company.Repository, userRepo user.Repository, subscriptionRepo subscription.Repository) *Service {
	return &Service{
		companyRepo:      companyRepo,
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
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

func (s *Service) AddUserToCompany(requestUserID, companyID uint, req *AddUserToCompanyRequest) error {
	// Check if requester has permission
	requesterRole, err := s.userRepo.FindUserRoleInCompany(requestUserID, companyID)
	if err != nil {
		return errors.NewForbiddenError("you don't have access to this company")
	}

	if requesterRole.Role != user.RoleOwner && requesterRole.Role != user.RoleAdmin {
		return errors.NewForbiddenError("only owners and admins can add users")
	}

	// Find user by email
	targetUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return errors.NewNotFoundError("user not found")
	}

	// Validate role
	role := user.Role(req.Role)
	if !role.IsValid() {
		return errors.NewValidationError("invalid role")
	}

	// Check if user already belongs to company
	existingRole, _ := s.userRepo.FindUserRoleInCompany(targetUser.ID, companyID)
	if existingRole != nil {
		return errors.NewConflictError("user already belongs to this company")
	}

	// Create user-company relationship
	ucr := &user.UserCompanyRole{
		UserID:    targetUser.ID,
		CompanyID: companyID,
		Role:      role,
		IsActive:  true,
	}

	if err := s.userRepo.CreateUserCompanyRole(ucr); err != nil {
		return errors.NewInternalError("failed to add user to company", err)
	}

	return nil
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
