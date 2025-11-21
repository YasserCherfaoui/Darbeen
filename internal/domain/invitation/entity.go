package invitation

import "time"

type Invitation struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"uniqueIndex;not null;size:64"`
	Email     string    `gorm:"not null;index"`
	CompanyID uint      `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	UserID    uint      `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Used      bool      `gorm:"default:false;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Invitation) TableName() string {
	return "invitations"
}

