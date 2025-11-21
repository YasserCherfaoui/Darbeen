package user

import (
	companyDomain "github.com/YasserCherfaoui/darween/internal/domain/company"
	franchiseDomain "github.com/YasserCherfaoui/darween/internal/domain/franchise"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	userRepo      user.Repository
	companyRepo   companyDomain.Repository
	franchiseRepo franchiseDomain.Repository
}

func NewService(userRepo user.Repository, companyRepo companyDomain.Repository, franchiseRepo franchiseDomain.Repository) *Service {
	return &Service{
		userRepo:      userRepo,
		companyRepo:   companyRepo,
		franchiseRepo: franchiseRepo,
	}
}

func (s *Service) GetUserByID(userID uint) (*UserResponse, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
	}, nil
}

func (s *Service) UpdateUser(userID uint, req *UpdateUserRequest) (*UserResponse, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	if req.FirstName != "" {
		u.FirstName = req.FirstName
	}
	if req.LastName != "" {
		u.LastName = req.LastName
	}

	if err := s.userRepo.Update(u); err != nil {
		return nil, errors.NewInternalError("failed to update user", err)
	}

	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
	}, nil
}

func (s *Service) GetUsersByCompanyID(companyID uint) ([]*UserWithRoleResponse, error) {
	users, err := s.userRepo.FindByCompanyID(companyID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch users", err)
	}

	roles, err := s.userRepo.FindCompanyUsersByCompanyID(companyID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch user roles", err)
	}

	// Create a map of user ID to role
	roleMap := make(map[uint]string)
	for _, r := range roles {
		roleMap[r.UserID] = r.Role.String()
	}

	var result []*UserWithRoleResponse
	for _, u := range users {
		result = append(result, &UserWithRoleResponse{
			UserResponse: UserResponse{
				ID:        u.ID,
				Email:     u.Email,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				IsActive:  u.IsActive,
			},
			Role: roleMap[u.ID],
		})
	}

	return result, nil
}

func (s *Service) ChangePassword(userID uint, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	// Find user
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	// Verify current password
	if !u.CheckPassword(req.CurrentPassword) {
		return nil, errors.NewUnauthorizedError("current password is incorrect")
	}

	// Check if new password is different from current password
	if u.CheckPassword(req.NewPassword) {
		return nil, errors.NewValidationError("new password must be different from current password")
	}

	// Hash new password
	if err := u.HashPassword(req.NewPassword); err != nil {
		return nil, errors.NewInternalError("failed to hash password", err)
	}

	// Update user
	if err := s.userRepo.Update(u); err != nil {
		return nil, errors.NewInternalError("failed to update password", err)
	}

	return &ChangePasswordResponse{
		Message: "Password has been changed successfully.",
	}, nil
}

func (s *Service) GetUserPortals(userID uint) (*UserPortalsResponse, error) {
	var portals []*PortalResponse

	// Fetch user's companies with roles
	companyRoles, err := s.userRepo.FindUserCompaniesByUserID(userID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch user companies", err)
	}

	for _, role := range companyRoles {
		company, err := s.companyRepo.FindByID(role.CompanyID)
		if err != nil {
			continue // Skip if company not found
		}

		portals = append(portals, &PortalResponse{
			Type:            "company",
			ID:              company.ID,
			Name:            company.Name,
			Code:            company.Code,
			Role:            role.Role.String(),
			ParentCompanyID: nil,
		})
	}

	// Fetch user's franchises with roles
	franchiseRoles, err := s.userRepo.FindUserFranchisesByUserID(userID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch user franchises", err)
	}

	for _, role := range franchiseRoles {
		franchise, err := s.franchiseRepo.FindByID(role.FranchiseID)
		if err != nil {
			continue // Skip if franchise not found
		}

		portals = append(portals, &PortalResponse{
			Type:            "franchise",
			ID:              franchise.ID,
			Name:            franchise.Name,
			Code:            franchise.Code,
			Role:            role.Role.String(),
			ParentCompanyID: &franchise.ParentCompanyID,
		})
	}

	return &UserPortalsResponse{
		Portals: portals,
	}, nil
}
