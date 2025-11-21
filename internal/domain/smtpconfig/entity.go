package smtpconfig

import "time"

type SMTPConfig struct {
	ID                    uint         `gorm:"primaryKey"`
	CompanyID             uint         `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Host                  string       `gorm:"not null"`
	User                  string       `gorm:"not null"`
	Password              string       `gorm:"not null"` // Encrypted
	Port                  int          `gorm:"not null"`
	FromName              string       `gorm:"default:''"`
	Security              SecurityType `gorm:"type:varchar(20);not null;default:'tls'"`
	SkipTLSVerify         bool         `gorm:"default:false"` // Skip TLS certificate verification
	RateLimit             int          `gorm:"not null;default:80"` // Emails per hour
	IsActive              bool         `gorm:"default:true"`
	IsDefault             bool         `gorm:"default:false"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (SMTPConfig) TableName() string {
	return "smtp_configs"
}

type SecurityType string

const (
	SecurityNone     SecurityType = "none"
	SecuritySSL      SecurityType = "ssl"
	SecurityTLS      SecurityType = "tls"
	SecuritySTARTTLS SecurityType = "starttls"
)

func (s SecurityType) String() string {
	return string(s)
}

func (s SecurityType) IsValid() bool {
	switch s {
	case SecurityNone, SecuritySSL, SecurityTLS, SecuritySTARTTLS:
		return true
	}
	return false
}

