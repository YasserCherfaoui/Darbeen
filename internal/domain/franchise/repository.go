package franchise

type Repository interface {
	// Franchise CRUD
	Create(franchise *Franchise) error
	FindByID(id uint) (*Franchise, error)
	FindByCode(code string) (*Franchise, error)
	FindByParentCompanyID(companyID uint) ([]*Franchise, error)
	Update(franchise *Franchise) error

	// Franchise Pricing
	CreatePricing(pricing *FranchisePricing) error
	FindPricing(franchiseID, variantID uint) (*FranchisePricing, error)
	FindAllPricingByFranchise(franchiseID uint) ([]*FranchisePricing, error)
	UpdatePricing(pricing *FranchisePricing) error
	DeletePricing(franchiseID, variantID uint) error
}




