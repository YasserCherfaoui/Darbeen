package user

import (
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	userRepo user.Repository
}

func NewService(userRepo user.Repository) *Service {
	return &Service{
		userRepo: userRepo,
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
