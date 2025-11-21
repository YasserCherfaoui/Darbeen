package email

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/company"
	"github.com/YasserCherfaoui/darween/internal/domain/emailqueue"
	"github.com/YasserCherfaoui/darween/internal/infrastructure/mailing"
)

type Service struct {
	mailingService *mailing.MailingService
	companyRepo    company.Repository
}

func NewService(mailingService *mailing.MailingService, companyRepo company.Repository) *Service {
	return &Service{
		mailingService: mailingService,
		companyRepo:    companyRepo,
	}
}

// SendPasswordResetEmail sends a password reset email
func (s *Service) SendPasswordResetEmail(req *SendPasswordResetEmailRequest) error {
	// Get company name
	companyName := "the company"
	if c, err := s.companyRepo.FindByID(req.CompanyID); err == nil {
		companyName = c.Name
		if companyName == "" {
			companyName = "the company"
		}
	}
	
	data := mailing.EmailTemplateData{
		CompanyName: companyName,
		UserName:    req.UserName,
		ResetToken:  req.ResetToken,
		ResetURL:    req.ResetURL,
	}
	
	htmlBody, _ := mailing.GeneratePasswordResetEmail(data)
	
	// Use HTML version
	return s.mailingService.QueueEmailWithType(
		req.CompanyID,
		nil, // Use default SMTP config
		[]string{req.UserEmail},
		"Password Reset Request",
		htmlBody,
		true, // IsHTML
		emailqueue.EmailTypePasswordReset,
		nil, // ScheduledAt
	)
}

// SendInvitationEmail sends an invitation email
func (s *Service) SendInvitationEmail(req *SendInvitationEmailRequest) error {
	data := mailing.EmailTemplateData{
		CompanyName:     req.CompanyName,
		UserName:        req.UserEmail, // Will use email as name if name not provided
		InviterName:     req.InviterName,
		InvitationToken: req.InvitationToken,
		InvitationURL:   req.InvitationURL,
	}
	
	htmlBody, plainBody := mailing.GenerateInvitationEmail(data)
	_ = plainBody // Plain text version available if needed
	
	// Use HTML version
	return s.mailingService.QueueEmailWithType(
		req.CompanyID,
		nil, // Use default SMTP config
		[]string{req.UserEmail},
		fmt.Sprintf("Invitation to Join %s", req.CompanyName),
		htmlBody,
		true, // IsHTML
		emailqueue.EmailTypeInvitation,
		nil, // ScheduledAt
	)
}

// SendStockAlertEmail sends a stock alert email
func (s *Service) SendStockAlertEmail(req *SendStockAlertEmailRequest) error {
	data := mailing.EmailTemplateData{
		CompanyName:    req.CompanyName,
		ProductName:    req.ProductName,
		VariantName:    req.VariantName,
		CurrentStock:   req.CurrentStock,
		Threshold:      req.Threshold,
		ProductDetails: req.ProductDetails,
	}
	
	htmlBody, plainBody := mailing.GenerateStockAlertEmail(data)
	_ = plainBody
	
	return s.mailingService.QueueEmailWithType(
		req.CompanyID,
		nil,
		req.To,
		fmt.Sprintf("Low Stock Alert: %s", req.ProductName),
		htmlBody,
		true,
		emailqueue.EmailTypeStockAlert,
		nil,
	)
}

// SendWarehouseBillEmail sends a warehouse bill email
func (s *Service) SendWarehouseBillEmail(req *SendWarehouseBillEmailRequest) error {
	data := mailing.EmailTemplateData{
		CompanyName: req.CompanyName,
		BillNumber:   req.BillNumber,
		BillType:     req.BillType,
		BillDate:     req.BillDate,
		BillItems:    req.BillItems,
		TotalAmount:  req.TotalAmount,
	}
	
	htmlBody, plainBody := mailing.GenerateWarehouseBillEmail(data)
	_ = plainBody
	
	billTypeLabel := "Entry Bill"
	if req.BillType == "exit" {
		billTypeLabel = "Exit Bill"
	}
	
	return s.mailingService.QueueEmailWithType(
		req.CompanyID,
		nil,
		req.To,
		fmt.Sprintf("Warehouse %s - %s", billTypeLabel, req.BillNumber),
		htmlBody,
		true,
		emailqueue.EmailTypeWarehouseBill,
		nil,
	)
}

// SendCustomEmail sends a custom email (for frontend composer)
func (s *Service) SendCustomEmail(req *SendCustomEmailRequest) error {
	return s.mailingService.QueueEmailWithType(
		req.CompanyID,
		nil,
		req.To,
		req.Subject,
		req.Body,
		req.IsHTML,
		emailqueue.EmailTypeCustom,
		nil,
	)
}

// SendNotificationEmail sends a generic notification email
func (s *Service) SendNotificationEmail(req *SendNotificationEmailRequest) error {
	htmlBody, plainBody := mailing.GenerateNotificationEmail(req.Subject, req.Message)
	_ = plainBody
	
	return s.mailingService.QueueEmailWithType(
		req.CompanyID,
		nil,
		req.To,
		req.Subject,
		htmlBody,
		true,
		emailqueue.EmailTypeNotification,
		nil,
	)
}

// SendCredentialsEmail sends credentials email to newly created user
func (s *Service) SendCredentialsEmail(req *SendCredentialsEmailRequest) error {
	data := mailing.EmailTemplateData{
		CompanyName: req.CompanyName,
		UserName:    req.UserEmail,
		CustomData: map[string]interface{}{
			"password":    req.Password,
			"inviterName": req.InviterName,
			"loginURL":    req.LoginURL,
		},
	}
	
	htmlBody, plainBody := mailing.GenerateCredentialsEmail(data)
	_ = plainBody
	
	return s.mailingService.QueueEmailWithType(
		req.CompanyID,
		nil, // Use default SMTP config
		[]string{req.UserEmail},
		fmt.Sprintf("Welcome to %s - Your Account Credentials", req.CompanyName),
		htmlBody,
		true, // IsHTML
		emailqueue.EmailTypeNotification,
		nil, // ScheduledAt
	)
}

// SendWelcomeEmail sends a welcome email to existing users who have been added to a company/franchise
func (s *Service) SendWelcomeEmail(req *SendWelcomeEmailRequest) error {
	data := mailing.EmailTemplateData{
		CompanyName: req.CompanyName,
		UserName:    req.UserEmail,
		InviterName: req.InviterName,
		Role:        req.Role,
		CustomData: map[string]interface{}{
			"loginURL": req.LoginURL,
		},
	}

	htmlBody, plainBody := mailing.GenerateWelcomeEmail(data)
	_ = plainBody

	return s.mailingService.QueueEmailWithType(
		req.CompanyID,
		nil, // Use default SMTP config
		[]string{req.UserEmail},
		fmt.Sprintf("Welcome to %s", req.CompanyName),
		htmlBody,
		true, // IsHTML
		emailqueue.EmailTypeNotification,
		nil, // ScheduledAt
	)
}

// SendNewUserSetupEmail sends a setup email to newly created users with credentials and OTP
func (s *Service) SendNewUserSetupEmail(req *SendNewUserSetupEmailRequest) error {
	data := mailing.EmailTemplateData{
		CompanyName: req.CompanyName,
		UserName:    req.UserEmail,
		OTPCode:     req.OTPCode,
		SetupURL:    req.SetupURL,
		CustomData: map[string]interface{}{
			"password":    req.Password,
			"inviterName": req.InviterName,
			"setupURL":    req.SetupURL,
		},
	}

	htmlBody, plainBody := mailing.GenerateNewUserSetupEmail(data)
	_ = plainBody

	return s.mailingService.QueueEmailWithType(
		req.CompanyID,
		nil, // Use default SMTP config
		[]string{req.UserEmail},
		fmt.Sprintf("Welcome to %s - Complete Your Account Setup", req.CompanyName),
		htmlBody,
		true, // IsHTML
		emailqueue.EmailTypeNotification,
		nil, // ScheduledAt
	)
}

