package company

import "time"

type Company struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null;default:''"` // Allow empty string as default
	Code        string `gorm:"not null;default:''"` // Remove unique constraint temporarily
	Description string
	IsActive    bool `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Company) TableName() string {
	return "companies"
}
