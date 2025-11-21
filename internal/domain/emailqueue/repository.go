package emailqueue

import "time"

type Repository interface {
	Create(email *EmailQueue) error
	FindByID(id uint) (*EmailQueue, error)
	FindPendingBySMTPConfig(smtpConfigID uint, limit int) ([]*EmailQueue, error)
	FindPendingByCompany(companyID uint, limit int) ([]*EmailQueue, error)
	FindCompaniesWithPendingEmails() ([]uint, error)
	Update(email *EmailQueue) error
	Delete(id uint) error
	GetEmailCountInWindow(smtpConfigID uint, windowStart time.Time) (int, error)
}

