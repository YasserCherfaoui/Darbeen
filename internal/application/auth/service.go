package auth

import (
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/internal/infrastructure/security"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	userRepo   user.Repository
	jwtManager *security.JWTManager
}

func NewService(userRepo user.Repository, jwtManager *security.JWTManager) *Service {
	return &Service{
		userRepo:   userRepo,
		jwtManager: jwtManager,
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
