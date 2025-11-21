package otp

type Repository interface {
	Create(otp *OTP) error
	FindByCodeAndEmail(code, email string) (*OTP, error)
	FindByCode(code string) (*OTP, error)
	MarkAsUsed(id uint) error
	DeleteExpired() error
}

