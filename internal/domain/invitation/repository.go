package invitation

type Repository interface {
	Create(invitation *Invitation) error
	FindByToken(token string) (*Invitation, error)
	FindByEmailAndCompany(email string, companyID uint) (*Invitation, error)
	MarkAsUsed(token string) error
	DeleteExpired() error
}

