package subscription

type Repository interface {
	Create(subscription *Subscription) error
	FindByID(id uint) (*Subscription, error)
	FindByCompanyID(companyID uint) (*Subscription, error)
	Update(subscription *Subscription) error
}

