package user

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
}

type UserWithRoleResponse struct {
	UserResponse
	Role string `json:"role"`
}

type PortalResponse struct {
	Type            string `json:"type"` // "company" or "franchise"
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Code            string `json:"code"`
	Role            string `json:"role"`
	ParentCompanyID *uint  `json:"parent_company_id,omitempty"` // Only for franchises
}

type UserPortalsResponse struct {
	Portals []*PortalResponse `json:"portals"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}

