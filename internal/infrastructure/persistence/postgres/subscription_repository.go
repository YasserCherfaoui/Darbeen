package postgres

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/subscription"
	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) subscription.Repository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(s *subscription.Subscription) error {
	return r.db.Create(s).Error
}

func (r *subscriptionRepository) FindByID(id uint) (*subscription.Subscription, error) {
	var s subscription.Subscription
	err := r.db.Where("id = ?", id).First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *subscriptionRepository) FindByCompanyID(companyID uint) (*subscription.Subscription, error) {
	var s subscription.Subscription
	err := r.db.Where("company_id = ?", companyID).First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *subscriptionRepository) Update(s *subscription.Subscription) error {
	return r.db.Save(s).Error
}
