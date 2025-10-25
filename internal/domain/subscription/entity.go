package subscription

import "time"

type Subscription struct {
	ID        uint      `gorm:"primaryKey"`
	CompanyID uint      `gorm:"uniqueIndex;not null;constraint:OnDelete:CASCADE"`
	PlanType  PlanType  `gorm:"not null"`
	Status    Status    `gorm:"not null"`
	StartDate time.Time `gorm:"not null"`
	EndDate   *time.Time
	MaxUsers  int `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Subscription) TableName() string {
	return "subscriptions"
}

type PlanType string

const (
	PlanFree       PlanType = "free"
	PlanBasic      PlanType = "basic"
	PlanPremium    PlanType = "premium"
	PlanEnterprise PlanType = "enterprise"
)

func (p PlanType) String() string {
	return string(p)
}

func (p PlanType) IsValid() bool {
	switch p {
	case PlanFree, PlanBasic, PlanPremium, PlanEnterprise:
		return true
	}
	return false
}

func (p PlanType) GetMaxUsers() int {
	switch p {
	case PlanFree:
		return 5
	case PlanBasic:
		return 20
	case PlanPremium:
		return 100
	case PlanEnterprise:
		return 1000
	}
	return 5
}

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusExpired  Status = "expired"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusExpired:
		return true
	}
	return false
}
