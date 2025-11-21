package emailqueue

import "time"

type EmailQueue struct {
	ID           uint        `gorm:"primaryKey"`
	CompanyID    uint        `gorm:"not null;index"`
	SMTPConfigID *uint       `gorm:"index"` // Nullable - uses default if null
	EmailType    EmailType   `gorm:"type:varchar(50);not null;default:'custom'"`
	To           string      `gorm:"type:text;not null"` // JSON array of email addresses
	Subject      string      `gorm:"not null"`
	Body         string      `gorm:"type:text;not null"`
	IsHTML       bool        `gorm:"default:false"`
	Status       EmailStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	Attempts     int         `gorm:"default:0"`
	MaxAttempts  int         `gorm:"default:3"`
	ErrorMessage string      `gorm:"type:text"`
	ScheduledAt  *time.Time  `gorm:"index"`
	SentAt       *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (EmailQueue) TableName() string {
	return "email_queue"
}

type EmailStatus string

const (
	EmailStatusPending    EmailStatus = "pending"
	EmailStatusProcessing EmailStatus = "processing"
	EmailStatusSent       EmailStatus = "sent"
	EmailStatusFailed     EmailStatus = "failed"
)

func (e EmailStatus) String() string {
	return string(e)
}

func (e EmailStatus) IsValid() bool {
	switch e {
	case EmailStatusPending, EmailStatusProcessing, EmailStatusSent, EmailStatusFailed:
		return true
	}
	return false
}

type EmailType string

const (
	EmailTypePasswordReset EmailType = "password_reset"
	EmailTypeInvitation    EmailType = "invitation"
	EmailTypeNotification  EmailType = "notification"
	EmailTypeStockAlert    EmailType = "stock_alert"
	EmailTypeWarehouseBill EmailType = "warehouse_bill"
	EmailTypeCustom        EmailType = "custom"
)

func (e EmailType) String() string {
	return string(e)
}

func (e EmailType) IsValid() bool {
	switch e {
	case EmailTypePasswordReset, EmailTypeInvitation, EmailTypeNotification, EmailTypeStockAlert, EmailTypeWarehouseBill, EmailTypeCustom:
		return true
	}
	return false
}
