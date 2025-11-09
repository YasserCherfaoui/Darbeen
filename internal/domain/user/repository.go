package user

type Repository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	FindByCompanyID(companyID uint) ([]*User, error)

	// User-Company-Role operations
	CreateUserCompanyRole(ucr *UserCompanyRole) error
	FindUserCompaniesByUserID(userID uint) ([]*UserCompanyRole, error)
	FindUserRoleInCompany(userID, companyID uint) (*UserCompanyRole, error)
	DeleteUserCompanyRole(userID, companyID uint) error
	FindCompanyUsersByCompanyID(companyID uint) ([]*UserCompanyRole, error)

	// User-Franchise-Role operations
	CreateUserFranchiseRole(ufr *UserFranchiseRole) error
	FindUserFranchisesByUserID(userID uint) ([]*UserFranchiseRole, error)
	FindUserRoleInFranchise(userID, franchiseID uint) (*UserFranchiseRole, error)
	DeleteUserFranchiseRole(userID, franchiseID uint) error
	FindFranchiseUsersByFranchiseID(franchiseID uint) ([]*UserFranchiseRole, error)
}
