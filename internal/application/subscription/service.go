package subscription

import (
	"github.com/YasserCherfaoui/darween/internal/domain/subscription"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	subscriptionRepo subscription.Repository
	userRepo         user.Repository
}

func NewService(subscriptionRepo subscription.Repository, userRepo user.Repository) *Service {
	return &Service{
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
	}
}

func (s *Service) GetSubscriptionByCompanyID(companyID uint) (*SubscriptionResponse, error) {
	sub, err := s.subscriptionRepo.FindByCompanyID(companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("subscription not found")
	}

	return &SubscriptionResponse{
		ID:        sub.ID,
		CompanyID: sub.CompanyID,
		PlanType:  sub.PlanType.String(),
		Status:    sub.Status.String(),
		StartDate: sub.StartDate,
		EndDate:   sub.EndDate,
		MaxUsers:  sub.MaxUsers,
	}, nil
}

func (s *Service) UpdateSubscription(userID, companyID uint, req *UpdateSubscriptionRequest) (*SubscriptionResponse, error) {
	// Check if user is owner
	role, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	if role.Role != user.RoleOwner {
		return nil, errors.NewForbiddenError("only company owners can update subscription")
	}

	// Get subscription
	sub, err := s.subscriptionRepo.FindByCompanyID(companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("subscription not found")
	}

	// Validate and update plan type
	planType := subscription.PlanType(req.PlanType)
	if !planType.IsValid() {
		return nil, errors.NewValidationError("invalid plan type")
	}
	sub.PlanType = planType
	sub.MaxUsers = planType.GetMaxUsers()

	// Update status if provided
	if req.Status != "" {
		status := subscription.Status(req.Status)
		if !status.IsValid() {
			return nil, errors.NewValidationError("invalid status")
		}
		sub.Status = status
	}

	if err := s.subscriptionRepo.Update(sub); err != nil {
		return nil, errors.NewInternalError("failed to update subscription", err)
	}

	return &SubscriptionResponse{
		ID:        sub.ID,
		CompanyID: sub.CompanyID,
		PlanType:  sub.PlanType.String(),
		Status:    sub.Status.String(),
		StartDate: sub.StartDate,
		EndDate:   sub.EndDate,
		MaxUsers:  sub.MaxUsers,
	}, nil
}
