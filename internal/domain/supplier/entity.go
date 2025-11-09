package supplier

import "time"

type Supplier struct {
	ID            uint   `gorm:"primaryKey"`
	CompanyID     uint   `gorm:"not null;index"`
	Name          string `gorm:"not null"`
	ContactPerson string
	Email         string
	Phone         string
	Address       string
	IsActive      bool `gorm:"default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (Supplier) TableName() string {
	return "suppliers"
}

// Business methods for Supplier
func (s *Supplier) IsValid() bool {
	return s.Name != "" && s.CompanyID > 0
}

