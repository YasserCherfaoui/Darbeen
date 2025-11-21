package postgres

import (
	"fmt"
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/emailqueue"
	"gorm.io/gorm"
)

type emailQueueRepository struct {
	db *gorm.DB
}

func NewEmailQueueRepository(db *gorm.DB) emailqueue.Repository {
	return &emailQueueRepository{db: db}
}

func (r *emailQueueRepository) Create(email *emailqueue.EmailQueue) error {
	return r.db.Create(email).Error
}

func (r *emailQueueRepository) FindByID(id uint) (*emailqueue.EmailQueue, error) {
	var email emailqueue.EmailQueue
	err := r.db.Where("id = ?", id).First(&email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("email queue not found")
		}
		return nil, err
	}
	return &email, nil
}

func (r *emailQueueRepository) FindPendingBySMTPConfig(smtpConfigID uint, limit int) ([]*emailqueue.EmailQueue, error) {
	var emails []*emailqueue.EmailQueue
	now := time.Now()
	query := r.db.Where("smtp_config_id = ? AND status = ?", smtpConfigID, emailqueue.EmailStatusPending).
		Where("(scheduled_at IS NULL OR scheduled_at <= ?)", now).
		Order("created_at ASC").
		Limit(limit)
	
	err := query.Find(&emails).Error
	return emails, err
}

func (r *emailQueueRepository) FindPendingByCompany(companyID uint, limit int) ([]*emailqueue.EmailQueue, error) {
	var emails []*emailqueue.EmailQueue
	now := time.Now()
	query := r.db.Where("company_id = ? AND status = ?", companyID, emailqueue.EmailStatusPending).
		Where("(scheduled_at IS NULL OR scheduled_at <= ?)", now).
		Order("created_at ASC").
		Limit(limit)
	
	err := query.Find(&emails).Error
	return emails, err
}

func (r *emailQueueRepository) Update(email *emailqueue.EmailQueue) error {
	return r.db.Save(email).Error
}

func (r *emailQueueRepository) Delete(id uint) error {
	return r.db.Delete(&emailqueue.EmailQueue{}, id).Error
}

func (r *emailQueueRepository) GetEmailCountInWindow(smtpConfigID uint, windowStart time.Time) (int, error) {
	var count int64
	err := r.db.Model(&emailqueue.EmailQueue{}).
		Where("smtp_config_id = ? AND status = ? AND sent_at >= ?", smtpConfigID, emailqueue.EmailStatusSent, windowStart).
		Count(&count).Error
	return int(count), err
}

