package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	FirstName string `gorm:"not null;default:''"` // Allow empty string as default
	LastName  string `gorm:"not null;default:''"` // Allow empty string as default
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}

// HashPassword hashes the user's password
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares the provided password with the user's hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// UserCompanyRole represents the many-to-many relationship between users and companies
type UserCompanyRole struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	CompanyID uint `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Role      Role `gorm:"not null"`
	IsActive  bool `gorm:"default:true"`
	CreatedAt time.Time
}

func (UserCompanyRole) TableName() string {
	return "user_company_roles"
}

// UserFranchiseRole represents the many-to-many relationship between users and franchises
type UserFranchiseRole struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	FranchiseID uint `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Role        Role `gorm:"not null"`
	IsActive    bool `gorm:"default:true"`
	CreatedAt   time.Time
}

func (UserFranchiseRole) TableName() string {
	return "user_franchise_roles"
}

type Role string

const (
	RoleOwner    Role = "owner"
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"
)

func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleAdmin, RoleManager, RoleEmployee:
		return true
	}
	return false
}
