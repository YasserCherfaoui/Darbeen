package email

type SendPasswordResetEmailRequest struct {
	CompanyID  uint
	UserEmail  string
	ResetToken string
	ResetURL   string
	UserName   string
	CompanyName string // Optional, will be fetched if not provided
}

type SendInvitationEmailRequest struct {
	CompanyID       uint
	UserEmail       string
	InviterName     string
	InvitationToken string
	InvitationURL   string
	CompanyName     string
}

type SendStockAlertEmailRequest struct {
	CompanyID     uint
	To            []string
	ProductName   string
	VariantName   string
	CurrentStock  int
	Threshold     int
	ProductDetails map[string]interface{}
	CompanyName   string
}

type SendWarehouseBillEmailRequest struct {
	CompanyID   uint
	To          []string
	BillNumber  string
	BillType    string
	BillDate    string
	BillItems   []map[string]interface{}
	TotalAmount float64
	CompanyName string
}

type SendCustomEmailRequest struct {
	CompanyID uint
	To        []string
	Subject   string
	Body      string
	IsHTML    bool
}

type SendNotificationEmailRequest struct {
	CompanyID uint
	To        []string
	Subject   string
	Message   string
}

type SendCredentialsEmailRequest struct {
	CompanyID   uint
	UserEmail   string
	Password    string
	InviterName string
	LoginURL    string
	CompanyName string
}

type SendWelcomeEmailRequest struct {
	CompanyID   uint
	UserEmail   string
	InviterName string
	Role        string
	LoginURL    string
	CompanyName string
}

type SendNewUserSetupEmailRequest struct {
	CompanyID   uint
	UserEmail   string
	Password    string
	InviterName string
	OTPCode     string
	SetupURL    string
	CompanyName string
}

