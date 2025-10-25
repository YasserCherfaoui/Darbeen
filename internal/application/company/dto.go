package company

type CreateCompanyRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

type UpdateCompanyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type CompanyResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type AddUserToCompanyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}


