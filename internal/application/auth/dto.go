package auth

import "fmt"

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetConfirmRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type PasswordResetResponse struct {
	Message string `json:"message"`
}

type ValidateInvitationRequest struct {
	Token string `json:"token" binding:"required"`
}

type ValidateInvitationResponse struct {
	Email     string `json:"email"`
	CompanyID uint   `json:"company_id"`
	Valid     bool   `json:"valid"`
}

type AcceptInvitationRequest struct {
	Token     string `json:"token" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
}

type ValidateOTPRequest struct {
	Code  interface{} `json:"code" binding:"required"`
	Email string       `json:"email" binding:"required,email"`
}

// GetCodeString converts the code to a string, handling both string and number types
func (r *ValidateOTPRequest) GetCodeString() string {
	switch v := r.Code.(type) {
	case string:
		return v
	case float64:
		// JSON numbers are unmarshaled as float64
		// Preserve leading zeros by formatting as 6-digit string
		return fmt.Sprintf("%06.0f", v)
	case int:
		return fmt.Sprintf("%06d", v)
	case int64:
		return fmt.Sprintf("%06d", v)
	default:
		// Try to convert to string and pad if needed
		str := fmt.Sprintf("%v", v)
		// If it's a number string, pad it
		if len(str) < 6 {
			return fmt.Sprintf("%06s", str)
		}
		return str
	}
}

type ValidateOTPResponse struct {
	Valid   bool   `json:"valid"`
	Purpose string `json:"purpose,omitempty"`
	UserID  uint   `json:"user_id,omitempty"`
}

type CompleteUserSetupRequest struct {
	Code      interface{} `json:"code" binding:"required"`
	Email     string      `json:"email" binding:"required,email"`
	FirstName string      `json:"first_name" binding:"required"`
	LastName  string      `json:"last_name" binding:"required"`
	Phone     string      `json:"phone,omitempty"`
	Password  string      `json:"password" binding:"required,min=6"`
}

// GetCodeString converts the code to a string, handling both string and number types
func (r *CompleteUserSetupRequest) GetCodeString() string {
	switch v := r.Code.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%06.0f", v)
	case int:
		return fmt.Sprintf("%06d", v)
	case int64:
		return fmt.Sprintf("%06d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type ChangePasswordWithOTPRequest struct {
	Code     interface{} `json:"code" binding:"required"`
	Email    string       `json:"email" binding:"required,email"`
	Password string       `json:"password" binding:"required,min=6"`
}

// GetCodeString converts the code to a string, handling both string and number types
func (r *ChangePasswordWithOTPRequest) GetCodeString() string {
	switch v := r.Code.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%06.0f", v)
	case int:
		return fmt.Sprintf("%06d", v)
	case int64:
		return fmt.Sprintf("%06d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

