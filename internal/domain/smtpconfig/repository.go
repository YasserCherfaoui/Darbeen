package smtpconfig

type Repository interface {
	Create(config *SMTPConfig) error
	FindByID(id uint) (*SMTPConfig, error)
	FindByIDAndCompany(id, companyID uint) (*SMTPConfig, error)
	FindByCompanyID(companyID uint) ([]*SMTPConfig, error)
	FindDefaultByCompanyID(companyID uint) (*SMTPConfig, error)
	Update(config *SMTPConfig) error
	Delete(id uint) error
	SetAsDefault(id, companyID uint) error
}

