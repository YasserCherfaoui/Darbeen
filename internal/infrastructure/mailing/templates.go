package mailing

import (
	"fmt"
	"html/template"
	"strings"
)

// EmailTemplateData holds data for email templates
type EmailTemplateData struct {
	// Common fields
	CompanyName string
	UserName    string
	
	// Password reset
	ResetToken string
	ResetURL   string
	
	// Invitation
	InviterName     string
	InvitationToken string
	InvitationURL   string
	
	// Stock alert
	ProductName    string
	VariantName    string
	CurrentStock   int
	Threshold      int
	ProductDetails map[string]interface{}
	
	// Warehouse bill
	BillNumber   string
	BillType    string
	BillDate    string
	BillItems   []map[string]interface{}
	TotalAmount float64
	
	// Custom data
	CustomData map[string]interface{}
}

// GeneratePasswordResetEmail generates HTML and plain text versions of password reset email
func GeneratePasswordResetEmail(data EmailTemplateData) (htmlBody, plainBody string) {
	resetURL := data.ResetURL
	if resetURL == "" && data.ResetToken != "" {
		// Default URL format if not provided
		resetURL = fmt.Sprintf("https://app.example.com/reset-password?token=%s", data.ResetToken)
	}
	
	htmlBody = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Password Reset</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #2c3e50;">Password Reset Request</h2>
		<p>Hello %s,</p>
		<p>We received a request to reset your password for your account at %s.</p>
		<p>Click the button below to reset your password:</p>
		<div style="text-align: center; margin: 30px 0;">
			<a href="%s" style="background-color: #3498db; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a>
		</div>
		<p>Or copy and paste this link into your browser:</p>
		<p style="word-break: break-all; color: #3498db;">%s</p>
		<p>This link will expire in 1 hour.</p>
		<p>If you didn't request a password reset, please ignore this email.</p>
		<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
		<p style="color: #7f8c8d; font-size: 12px;">This is an automated message, please do not reply.</p>
	</div>
</body>
</html>
`, data.UserName, data.CompanyName, resetURL, resetURL)
	
	plainBody = fmt.Sprintf(`
Password Reset Request

Hello %s,

We received a request to reset your password for your account at %s.

Click the following link to reset your password:
%s

This link will expire in 1 hour.

If you didn't request a password reset, please ignore this email.

---
This is an automated message, please do not reply.
`, data.UserName, data.CompanyName, resetURL)
	
	return htmlBody, plainBody
}

// GenerateInvitationEmail generates HTML and plain text versions of invitation email
func GenerateInvitationEmail(data EmailTemplateData) (htmlBody, plainBody string) {
	invitationURL := data.InvitationURL
	if invitationURL == "" && data.InvitationToken != "" {
		invitationURL = fmt.Sprintf("https://app.example.com/accept-invitation?token=%s", data.InvitationToken)
	}
	
	htmlBody = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Invitation to Join %s</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #2c3e50;">You're Invited!</h2>
		<p>Hello %s,</p>
		<p>%s has invited you to join %s.</p>
		<p>Click the button below to accept the invitation and create your account:</p>
		<div style="text-align: center; margin: 30px 0;">
			<a href="%s" style="background-color: #27ae60; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">Accept Invitation</a>
		</div>
		<p>Or copy and paste this link into your browser:</p>
		<p style="word-break: break-all; color: #27ae60;">%s</p>
		<p>This invitation will expire in 7 days.</p>
		<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
		<p style="color: #7f8c8d; font-size: 12px;">This is an automated message, please do not reply.</p>
	</div>
</body>
</html>
`, data.CompanyName, data.UserName, data.InviterName, data.CompanyName, invitationURL, invitationURL)
	
	plainBody = fmt.Sprintf(`
You're Invited!

Hello %s,

%s has invited you to join %s.

Click the following link to accept the invitation and create your account:
%s

This invitation will expire in 7 days.

---
This is an automated message, please do not reply.
`, data.UserName, data.InviterName, data.CompanyName, invitationURL)
	
	return htmlBody, plainBody
}

// GenerateStockAlertEmail generates HTML and plain text versions of stock alert email
func GenerateStockAlertEmail(data EmailTemplateData) (htmlBody, plainBody string) {
	variantInfo := ""
	if data.VariantName != "" {
		variantInfo = fmt.Sprintf(" (%s)", data.VariantName)
	}
	
	htmlBody = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Stock Alert - %s</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #e74c3c;">⚠️ Low Stock Alert</h2>
		<p>Hello,</p>
		<p>This is an automated alert from %s regarding low stock levels.</p>
		<div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0;">
			<h3 style="margin-top: 0; color: #856404;">Product: %s%s</h3>
			<p><strong>Current Stock:</strong> %d units</p>
			<p><strong>Threshold:</strong> %d units</p>
			<p style="color: #e74c3c; font-weight: bold;">⚠️ Stock is below the minimum threshold!</p>
		</div>
		<p>Please consider restocking this product to avoid stockouts.</p>
		<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
		<p style="color: #7f8c8d; font-size: 12px;">This is an automated message, please do not reply.</p>
	</div>
</body>
</html>
`, data.ProductName, data.CompanyName, data.ProductName, variantInfo, data.CurrentStock, data.Threshold)
	
	plainBody = fmt.Sprintf(`
Low Stock Alert

Hello,

This is an automated alert from %s regarding low stock levels.

Product: %s%s
Current Stock: %d units
Threshold: %d units

⚠️ Stock is below the minimum threshold!

Please consider restocking this product to avoid stockouts.

---
This is an automated message, please do not reply.
`, data.CompanyName, data.ProductName, variantInfo, data.CurrentStock, data.Threshold)
	
	return htmlBody, plainBody
}

// GenerateWarehouseBillEmail generates HTML and plain text versions of warehouse bill email
func GenerateWarehouseBillEmail(data EmailTemplateData) (htmlBody, plainBody string) {
	billTypeLabel := "Entry Bill"
	if data.BillType == "exit" {
		billTypeLabel = "Exit Bill"
	}
	
	itemsHTML := ""
	itemsPlain := ""
	if len(data.BillItems) > 0 {
		itemsHTML = "<table style='width: 100%%; border-collapse: collapse; margin: 20px 0;'><thead><tr style='background-color: #f8f9fa;'><th style='padding: 10px; text-align: left; border-bottom: 2px solid #dee2e6;'>Product</th><th style='padding: 10px; text-align: right; border-bottom: 2px solid #dee2e6;'>Quantity</th></tr></thead><tbody>"
		for _, item := range data.BillItems {
			productName := ""
			if name, ok := item["product_name"].(string); ok {
				productName = name
			}
			quantity := 0
			if qty, ok := item["quantity"].(int); ok {
				quantity = qty
			} else if qty, ok := item["quantity"].(float64); ok {
				quantity = int(qty)
			}
			itemsHTML += fmt.Sprintf("<tr><td style='padding: 10px; border-bottom: 1px solid #dee2e6;'>%s</td><td style='padding: 10px; text-align: right; border-bottom: 1px solid #dee2e6;'>%d</td></tr>", productName, quantity)
			itemsPlain += fmt.Sprintf("  - %s: %d\n", productName, quantity)
		}
		itemsHTML += "</tbody></table>"
	}
	
	htmlBody = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Warehouse Bill - %s</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #2c3e50;">Warehouse %s</h2>
		<p>Hello,</p>
		<p>A new warehouse %s has been processed for %s.</p>
		<div style="background-color: #f8f9fa; padding: 15px; margin: 20px 0; border-radius: 5px;">
			<p><strong>Bill Number:</strong> %s</p>
			<p><strong>Date:</strong> %s</p>
			<p><strong>Type:</strong> %s</p>
		</div>
		<h3>Items:</h3>
		%s
		<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
		<p style="color: #7f8c8d; font-size: 12px;">This is an automated message, please do not reply.</p>
	</div>
</body>
</html>
`, data.BillNumber, billTypeLabel, billTypeLabel, data.CompanyName, data.BillNumber, data.BillDate, billTypeLabel, itemsHTML)
	
	plainBody = fmt.Sprintf(`
Warehouse %s

Hello,

A new warehouse %s has been processed for %s.

Bill Number: %s
Date: %s
Type: %s

Items:
%s
---
This is an automated message, please do not reply.
`, billTypeLabel, billTypeLabel, data.CompanyName, data.BillNumber, data.BillDate, billTypeLabel, itemsPlain)
	
	return htmlBody, plainBody
}

// GenerateNotificationEmail generates HTML and plain text versions of generic notification email
func GenerateNotificationEmail(subject, message string) (htmlBody, plainBody string) {
	htmlBody = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #2c3e50;">%s</h2>
		<div style="margin: 20px 0;">
			%s
		</div>
		<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
		<p style="color: #7f8c8d; font-size: 12px;">This is an automated message, please do not reply.</p>
	</div>
</body>
</html>
`, template.HTMLEscapeString(subject), template.HTMLEscapeString(subject), message)
	
	plainBody = fmt.Sprintf(`
%s

%s

---
This is an automated message, please do not reply.
`, subject, strings.ReplaceAll(message, "<br>", "\n"))
	
	// Remove HTML tags from plain text
	plainBody = strings.ReplaceAll(plainBody, "<p>", "")
	plainBody = strings.ReplaceAll(plainBody, "</p>", "\n")
	plainBody = strings.ReplaceAll(plainBody, "<strong>", "")
	plainBody = strings.ReplaceAll(plainBody, "</strong>", "")
	plainBody = strings.ReplaceAll(plainBody, "<em>", "")
	plainBody = strings.ReplaceAll(plainBody, "</em>", "")
	
	return htmlBody, plainBody
}

