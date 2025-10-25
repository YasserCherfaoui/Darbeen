package postgres

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *userRepository) FindByID(id uint) (*user.User, error) {
	var u user.User
	err := r.db.Where("id = ?", id).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(u *user.User) error {
	return r.db.Save(u).Error
}

func (r *userRepository) FindByCompanyID(companyID uint) ([]*user.User, error) {
	var users []*user.User
	err := r.db.
		Joins("JOIN user_company_roles ON user_company_roles.user_id = users.id").
		Where("user_company_roles.company_id = ? AND user_company_roles.is_active = ?", companyID, true).
		Find(&users).Error
	return users, err
}

func (r *userRepository) CreateUserCompanyRole(ucr *user.UserCompanyRole) error {
	return r.db.Create(ucr).Error
}

func (r *userRepository) FindUserCompaniesByUserID(userID uint) ([]*user.UserCompanyRole, error) {
	var roles []*user.UserCompanyRole
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&roles).Error
	return roles, err
}

func (r *userRepository) FindUserRoleInCompany(userID, companyID uint) (*user.UserCompanyRole, error) {
	var role user.UserCompanyRole
	err := r.db.Where("user_id = ? AND company_id = ? AND is_active = ?", userID, companyID, true).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *userRepository) DeleteUserCompanyRole(userID, companyID uint) error {
	return r.db.Model(&user.UserCompanyRole{}).
		Where("user_id = ? AND company_id = ?", userID, companyID).
		Update("is_active", false).Error
}

func (r *userRepository) FindCompanyUsersByCompanyID(companyID uint) ([]*user.UserCompanyRole, error) {
	var roles []*user.UserCompanyRole
	err := r.db.Where("company_id = ? AND is_active = ?", companyID, true).Find(&roles).Error
	return roles, err
}
