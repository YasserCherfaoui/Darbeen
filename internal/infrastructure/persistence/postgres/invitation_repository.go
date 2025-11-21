package postgres

import (
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/invitation"
	"gorm.io/gorm"
)

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) invitation.Repository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) Create(inv *invitation.Invitation) error {
	return r.db.Create(inv).Error
}

func (r *invitationRepository) FindByToken(token string) (*invitation.Invitation, error) {
	var inv invitation.Invitation
	err := r.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *invitationRepository) FindByEmailAndCompany(email string, companyID uint) (*invitation.Invitation, error) {
	var inv invitation.Invitation
	err := r.db.Where("email = ? AND company_id = ? AND used = ? AND expires_at > ?", email, companyID, false, time.Now()).
		Order("created_at DESC").
		First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *invitationRepository) MarkAsUsed(token string) error {
	return r.db.Model(&invitation.Invitation{}).
		Where("token = ?", token).
		Update("used", true).Error
}

func (r *invitationRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ? OR used = ?", time.Now(), true).
		Delete(&invitation.Invitation{}).Error
}

