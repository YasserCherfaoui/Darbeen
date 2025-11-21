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
	// In a production system, you'd want to:
	// 1. Get all companies with pending emails
	// 2. For each company, process emails using ProcessEmailsForCompany
	// 3. Handle errors gracefully

	// For now, this is a simplified version that processes emails
	// The actual implementation would query for companies with pending emails
	// and process them one by one
	
	log.Println("Processing email queue...")
	// TODO: Implement full queue processing logic
	// This would involve:
	// - Querying for companies with pending emails
	// - Calling ProcessEmailsForCompany for each company
	// - Handling errors and retries
}

