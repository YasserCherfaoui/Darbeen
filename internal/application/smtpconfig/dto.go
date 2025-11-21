package smtpconfig

type CreateSMTPConfigRequest struct {
	Host          string `json:"host" binding:"required"`
	User          string `json:"user" binding:"required"`
	Password      string `json:"password" binding:"required"`
	Port          int    `json:"port" binding:"required,min=1,max=65535"`
	FromName      string `json:"from_name"`
	Security      string `json:"security" binding:"required,oneof=none ssl tls starttls"`
	SkipTLSVerify *bool  `json:"skip_tls_verify"` // Skip TLS certificate verification
	RateLimit     int    `json:"rate_limit" binding:"min=1"`
	IsActive      *bool  `json:"is_active"`
}

type UpdateSMTPConfigRequest struct {
	Host          *string `json:"host"`
	User          *string `json:"user"`
	Password      *string `json:"password"`
	Port          *int    `json:"port" binding:"omitempty,min=1,max=65535"`
	FromName      *string `json:"from_name"`
	Security      *string `json:"security" binding:"omitempty,oneof=none ssl tls starttls"`
	SkipTLSVerify *bool   `json:"skip_tls_verify"` // Skip TLS certificate verification
	RateLimit     *int    `json:"rate_limit" binding:"omitempty,min=1"`
	IsActive      *bool   `json:"is_active"`
}

type SMTPConfigResponse struct {
	ID            uint   `json:"id"`
	CompanyID     uint   `json:"company_id"`
	Host          string `json:"host"`
	User          string `json:"user"`
	Port          int    `json:"port"`
	FromName      string `json:"from_name"`
	Security      string `json:"security"`
	SkipTLSVerify bool   `json:"skip_tls_verify"`
	RateLimit     int    `json:"rate_limit"`
	IsActive      bool   `json:"is_active"`
	IsDefault     bool   `json:"is_default"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ListSMTPConfigsResponse struct {
	Configs []*SMTPConfigResponse `json:"configs"`
}

