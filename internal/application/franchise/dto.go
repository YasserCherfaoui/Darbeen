package franchise

// Request DTOs

type CreateFranchiseRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description,omitempty"`
	Address     string `json:"address,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
}

type UpdateFranchiseRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Address     string `json:"address,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type SetFranchisePricingRequest struct {
	ProductVariantID uint     `json:"product_variant_id" binding:"required"`
	RetailPrice      *float64 `json:"retail_price,omitempty"`
	WholesalePrice   *float64 `json:"wholesale_price,omitempty"`
}

type InitializeFranchiseInventoryRequest struct {
	// Empty for now - initializes with zero stock
}

type AddUserToFranchiseRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

// Response DTOs

type FranchiseResponse struct {
	ID              uint   `json:"id"`
	ParentCompanyID uint   `json:"parent_company_id"`
	Name            string `json:"name"`
	Code            string `json:"code"`
	Description     string `json:"description,omitempty"`
	Address         string `json:"address,omitempty"`
	Phone           string `json:"phone,omitempty"`
	Email           string `json:"email,omitempty"`
	IsActive        bool   `json:"is_active"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type FranchisePricingResponse struct {
	ID                    uint     `json:"id"`
	FranchiseID           uint     `json:"franchise_id"`
	ProductVariantID      uint     `json:"product_variant_id"`
	VariantName           string   `json:"variant_name,omitempty"`
	VariantSKU            string   `json:"variant_sku,omitempty"`
	RetailPrice           *float64 `json:"retail_price,omitempty"`
	WholesalePrice        *float64 `json:"wholesale_price,omitempty"`
	DefaultRetailPrice    float64  `json:"default_retail_price"`
	DefaultWholesalePrice float64  `json:"default_wholesale_price"`
	IsActive              bool     `json:"is_active"`
}

type FranchiseListResponse struct {
	Franchises []*FranchiseResponse `json:"franchises"`
	Total      int                  `json:"total"`
}




