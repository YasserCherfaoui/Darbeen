package company

type Repository interface {
	Create(company *Company) error
	FindByID(id uint) (*Company, error)
	FindByCode(code string) (*Company, error)
	Update(company *Company) error
	FindByUserID(userID uint) ([]*Company, error)
}

