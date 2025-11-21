package mailing

import (
	"context"
	"log"
	"time"
)

// EmailQueueWorker processes the email queue periodically
type EmailQueueWorker struct {
	mailingService *MailingService
	interval       time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewEmailQueueWorker creates a new email queue worker
func NewEmailQueueWorker(mailingService *MailingService, interval time.Duration) *EmailQueueWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &EmailQueueWorker{
		mailingService: mailingService,
		interval:       interval,
		ctx:            ctx,
		cancel:         cancel,
	}
}

// Start starts the worker
func (w *EmailQueueWorker) Start() {
	go w.run()
}

// Stop stops the worker
func (w *EmailQueueWorker) Stop() {
	w.cancel()
}

func (w *EmailQueueWorker) run() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Println("Email queue worker started")

	for {
		select {
		case <-w.ctx.Done():
			log.Println("Email queue worker stopped")
			return
		case <-ticker.C:
			w.processQueue()
		}
	}
}

func (w *EmailQueueWorker) processQueue() {
	log.Println("Processing email queue...")

	// Get all companies with pending emails
	companyIDs, err := w.mailingService.emailQueueRepo.FindCompaniesWithPendingEmails()
	if err != nil {
		log.Printf("Failed to find companies with pending emails: %v", err)
		return
	}

	if len(companyIDs) == 0 {
		log.Println("No companies with pending emails")
		return
	}

	log.Printf("Found %d companies with pending emails", len(companyIDs))

	// Process emails for each company
	for _, companyID := range companyIDs {
		if err := w.mailingService.ProcessEmailsForCompany(companyID); err != nil {
			log.Printf("Failed to process emails for company %d: %v", companyID, err)
			// Continue processing other companies even if one fails
			continue
		}
	}

	log.Println("Email queue processing completed")
}

