package company

type CreateCompanyRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	ERPUrl      string `json:"erp_url"`
}

type UpdateCompanyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ERPUrl      string `json:"erp_url"`
	IsActive    *bool  `json:"is_active"`
}

type CompanyResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	ERPUrl      string `json:"erp_url"`
	IsActive    bool   `json:"is_active"`
}

type AddUserToCompanyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddUserToCompanyResponse struct {
	UserCreated bool             `json:"user_created"`
	EmailSent   bool             `json:"email_sent"`
	Credentials *UserCredentials `json:"credentials,omitempty"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type UserWithRoleResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
}

type ListCompanyUsersResponse struct {
	Users []*UserWithRoleResponse `json:"users"`
}


