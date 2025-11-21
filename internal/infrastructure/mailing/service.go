package mailing

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/emailqueue"
	"github.com/YasserCherfaoui/darween/internal/domain/smtpconfig"
	"github.com/YasserCherfaoui/darween/pkg/encryption"
)

type MailingService struct {
	smtpRepo       smtpconfig.Repository
	emailQueueRepo emailqueue.Repository
}

func NewMailingService(smtpRepo smtpconfig.Repository, emailQueueRepo emailqueue.Repository) *MailingService {
	return &MailingService{
		smtpRepo:       smtpRepo,
		emailQueueRepo: emailQueueRepo,
	}
}

// QueueEmail adds an email to the queue for processing
func (m *MailingService) QueueEmail(companyID uint, smtpConfigID *uint, to []string, subject, body string, isHTML bool, scheduledAt *time.Time) error {
	return m.QueueEmailWithType(companyID, smtpConfigID, to, subject, body, isHTML, emailqueue.EmailTypeCustom, scheduledAt)
}

// QueueEmailWithType adds an email to the queue for processing with a specific email type
func (m *MailingService) QueueEmailWithType(companyID uint, smtpConfigID *uint, to []string, subject, body string, isHTML bool, emailType emailqueue.EmailType, scheduledAt *time.Time) error {
	// Convert to array to JSON string
	toJSON, err := json.Marshal(to)
	if err != nil {
		return fmt.Errorf("failed to marshal recipients: %w", err)
	}

	email := &emailqueue.EmailQueue{
		CompanyID:    companyID,
		SMTPConfigID: smtpConfigID,
		EmailType:    emailType,
		To:           string(toJSON),
		Subject:      subject,
		Body:         body,
		IsHTML:       isHTML,
		Status:       emailqueue.EmailStatusPending,
		Attempts:     0,
		MaxAttempts:  3,
		ScheduledAt:  scheduledAt,
	}

	return m.emailQueueRepo.Create(email)
}

// SendEmail queues an email using the default SMTP config for the company
func (m *MailingService) SendEmail(companyID uint, to []string, subject, body string, isHTML bool) error {
	return m.QueueEmail(companyID, nil, to, subject, body, isHTML, nil)
}

// SendEmailWithConfig queues an email using a specific SMTP config
func (m *MailingService) SendEmailWithConfig(configID, companyID uint, to []string, subject, body string, isHTML bool) error {
	return m.QueueEmail(companyID, &configID, to, subject, body, isHTML, nil)
}

// ProcessEmailQueue processes pending emails respecting rate limits
func (m *MailingService) ProcessEmailQueue() error {
	// Get all companies with pending emails (group by company)
	// For simplicity, we'll process by company
	// In production, you might want to batch by SMTP config

	// Process in batches of 10 emails at a time
	batchSize := 10
	processed := 0

	for {
		// Find pending emails (we'll need to add a method to find pending emails across all companies)
		// For now, let's process by finding emails that need to be sent
		// This is a simplified version - in production you'd want better batching

		// Get pending emails (limit to batch size)
		// Note: This requires a method to find pending emails without company filter
		// For now, we'll process emails that are ready to be sent

		// Break if we've processed enough
		if processed >= batchSize {
			break
		}

		// In a real implementation, you'd query for pending emails across all companies
		// and process them respecting rate limits per SMTP config
		break
	}

	return nil
}

// ProcessEmailsForCompany processes pending emails for a company (using default SMTP config if no config specified)
func (m *MailingService) ProcessEmailsForCompany(companyID uint) error {
	// Get pending emails for this company
	pendingEmails, err := m.emailQueueRepo.FindPendingByCompany(companyID, 100)
	if err != nil {
		return fmt.Errorf("failed to find pending emails: %w", err)
	}

	if len(pendingEmails) == 0 {
		return nil
	}

	// Group emails by SMTP config ID
	emailsByConfig := make(map[uint][]*emailqueue.EmailQueue)
	var emailsWithoutConfig []*emailqueue.EmailQueue

	for _, email := range pendingEmails {
		if email.SMTPConfigID == nil {
			emailsWithoutConfig = append(emailsWithoutConfig, email)
		} else {
			emailsByConfig[*email.SMTPConfigID] = append(emailsByConfig[*email.SMTPConfigID], email)
		}
	}

	// Process emails without config using default
	if len(emailsWithoutConfig) > 0 {
		defaultConfig, err := m.smtpRepo.FindDefaultByCompanyID(companyID)
		if err != nil {
			log.Printf("No default SMTP config found for company %d: %v", companyID, err)
			return nil
		}
		if err := m.processEmailsWithConfig(emailsWithoutConfig, defaultConfig); err != nil {
			log.Printf("Failed to process emails with default config: %v", err)
		}
	}

	// Process emails with specific configs
	for configID, emails := range emailsByConfig {
		config, err := m.smtpRepo.FindByID(configID)
		if err != nil {
			log.Printf("Failed to find SMTP config %d: %v", configID, err)
			continue
		}
		if err := m.processEmailsWithConfig(emails, config); err != nil {
			log.Printf("Failed to process emails with config %d: %v", configID, err)
		}
	}

	return nil
}

// processEmailsWithConfig processes emails using a specific SMTP config
func (m *MailingService) processEmailsWithConfig(emails []*emailqueue.EmailQueue, config *smtpconfig.SMTPConfig) error {
	if !config.IsActive {
		return fmt.Errorf("SMTP config is not active")
	}

	// Check rate limit
	windowStart := time.Now().Add(-1 * time.Hour)
	emailCount, err := m.emailQueueRepo.GetEmailCountInWindow(config.ID, windowStart)
	if err != nil {
		return fmt.Errorf("failed to get email count: %w", err)
	}

	if emailCount >= config.RateLimit {
		log.Printf("Rate limit reached for SMTP config %d: %d/%d emails in last hour", config.ID, emailCount, config.RateLimit)
		return nil // Don't process, wait for next cycle
	}

	// Calculate how many emails we can send
	remainingQuota := config.RateLimit - emailCount
	if remainingQuota <= 0 {
		return nil
	}

	// Limit to remaining quota
	if len(emails) > remainingQuota {
		emails = emails[:remainingQuota]
	}

	// Process each email
	for _, email := range emails {
		// Update status to processing
		email.Status = emailqueue.EmailStatusProcessing
		if err := m.emailQueueRepo.Update(email); err != nil {
			log.Printf("Failed to update email status to processing: %v", err)
			continue
		}

		// Send email
		err := m.sendEmailDirectly(email, config)
		if err != nil {
			// Mark as failed and increment attempts
			email.Status = emailqueue.EmailStatusFailed
			email.Attempts++
			email.ErrorMessage = err.Error()
			if email.Attempts >= email.MaxAttempts {
				// Max attempts reached, mark as failed permanently
				m.emailQueueRepo.Update(email)
			} else {
				// Will retry later
				email.Status = emailqueue.EmailStatusPending
				m.emailQueueRepo.Update(email)
			}
			log.Printf("Failed to send email %d: %v", email.ID, err)
			continue
		}

		// Mark as sent
		now := time.Now()
		email.Status = emailqueue.EmailStatusSent
		email.SentAt = &now
		if err := m.emailQueueRepo.Update(email); err != nil {
			log.Printf("Failed to update email status to sent: %v", err)
		}
	}

	return nil
}

// ProcessEmailsForSMTPConfig processes pending emails for a specific SMTP config
func (m *MailingService) ProcessEmailsForSMTPConfig(smtpConfigID uint) error {
	config, err := m.smtpRepo.FindByID(smtpConfigID)
	if err != nil {
		return fmt.Errorf("failed to find SMTP config: %w", err)
	}

	// Get pending emails for this config
	pendingEmails, err := m.emailQueueRepo.FindPendingBySMTPConfig(smtpConfigID, 100)
	if err != nil {
		return fmt.Errorf("failed to find pending emails: %w", err)
	}

	return m.processEmailsWithConfig(pendingEmails, config)
}

// sendEmailDirectly sends an email using the SMTP config
func (m *MailingService) sendEmailDirectly(email *emailqueue.EmailQueue, config *smtpconfig.SMTPConfig) error {
	// Decrypt password
	password, err := m.decryptPassword(config.Password)
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Parse recipients
	var to []string
	if err := json.Unmarshal([]byte(email.To), &to); err != nil {
		return fmt.Errorf("failed to unmarshal recipients: %w", err)
	}

	// Build email message
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", config.FromName, config.User))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to[0]))
	if len(to) > 1 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", to[1:]))
	}
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	if email.IsHTML {
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}
	msg.WriteString("\r\n")
	msg.WriteString(email.Body)

	// Connect to SMTP server based on security type
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	hostname := config.Host // Use hostname without port for TLS verification

	auth := smtp.PlainAuth("", config.User, password, config.Host)

	switch config.Security {
	case smtpconfig.SecuritySSL:
		// SSL connection
		return m.sendWithSSL(addr, hostname, config.SkipTLSVerify, auth, config.User, to, msg.Bytes())
	case smtpconfig.SecurityTLS:
		// TLS connection
		return m.sendWithTLS(addr, hostname, config.SkipTLSVerify, auth, config.User, to, msg.Bytes())
	case smtpconfig.SecuritySTARTTLS:
		// STARTTLS
		return m.sendWithSTARTTLS(addr, hostname, config.SkipTLSVerify, auth, config.User, to, msg.Bytes())
	case smtpconfig.SecurityNone:
		// No encryption
		return smtp.SendMail(addr, auth, config.User, to, msg.Bytes())
	default:
		return fmt.Errorf("unsupported security type: %s", config.Security)
	}
}

func (m *MailingService) sendWithSSL(addr string, hostname string, skipVerify bool, auth smtp.Auth, from string, to []string, msg []byte) error {
	// SSL requires TLS from the start
	// Use hostname (without port) for ServerName to match certificate
	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipVerify,
		ServerName:         hostname,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect with SSL: %w", err)
	}
	defer conn.Close()

	// When using an existing TLS connection, pass hostname (without port) to NewClient
	client, err := smtp.NewClient(conn, hostname)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data: %w", err)
	}

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close data: %w", err)
	}

	return client.Quit()
}

func (m *MailingService) sendWithTLS(addr string, hostname string, skipVerify bool, auth smtp.Auth, from string, to []string, msg []byte) error {
	// TLS connection
	// Use hostname (without port) for ServerName to match certificate
	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipVerify,
		ServerName:         hostname,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect with TLS: %w", err)
	}
	defer conn.Close()

	// When using an existing TLS connection, pass hostname (without port) to NewClient
	client, err := smtp.NewClient(conn, hostname)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data: %w", err)
	}

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close data: %w", err)
	}

	return client.Quit()
}

func (m *MailingService) sendWithSTARTTLS(addr string, hostname string, skipVerify bool, auth smtp.Auth, from string, to []string, msg []byte) error {
	// STARTTLS - connect first, then upgrade
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	// Use hostname (without port) for ServerName to match certificate
	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipVerify,
		ServerName:         hostname,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data: %w", err)
	}

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close data: %w", err)
	}

	return client.Quit()
}

// decryptPassword decrypts the encrypted SMTP password
func (m *MailingService) decryptPassword(encryptedPassword string) (string, error) {
	return encryption.Decrypt(encryptedPassword)
}
