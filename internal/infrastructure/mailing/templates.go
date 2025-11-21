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
	BillNumber  string
	BillType    string
	BillDate    string
	BillItems   []map[string]interface{}
	TotalAmount float64

	// Custom data
	CustomData map[string]interface{}

	// OTP and setup
	Role        string
	OTPCode     string
	SetupURL    string
	PhoneNumber string
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
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>Invitation to Join %s</title>
	<!--[if mso]>
	<style type="text/css">
		table {border-collapse:collapse;border-spacing:0;margin:0;}
		div, td {padding:0;}
		div {margin:0 !important;}
	</style>
	<![endif]-->
</head>
<body style="margin: 0; padding: 0; background-color: #f4f4f4; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;">
	<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin: 0; padding: 0; background-color: #f4f4f4;">
		<tr>
			<td align="center" style="padding: 40px 20px;">
				<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="max-width: 600px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
					<!-- Header -->
					<tr>
						<td style="padding: 40px 40px 20px; text-align: center; background-color: #ffffff; border-radius: 8px 8px 0 0;">
							<h1 style="margin: 0; color: #2c3e50; font-size: 28px; font-weight: 600;">You're Invited!</h1>
						</td>
					</tr>
					<!-- Content -->
					<tr>
						<td style="padding: 20px 40px;">
							<p style="margin: 0 0 16px; color: #333333; font-size: 16px; line-height: 1.6;">Hello,</p>
							<p style="margin: 0 0 16px; color: #333333; font-size: 16px; line-height: 1.6;">
								<strong>%s</strong> has invited you to join <strong>%s</strong>.
							</p>
							<p style="margin: 0 0 24px; color: #333333; font-size: 16px; line-height: 1.6;">
								Click the button below to accept the invitation and set up your account:
							</p>
							<!-- CTA Button -->
							<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%">
								<tr>
									<td align="center" style="padding: 0 0 24px;">
										<a href="%s" style="display: inline-block; background-color: #27ae60; color: #ffffff; text-decoration: none; padding: 14px 32px; border-radius: 6px; font-size: 16px; font-weight: 600; text-align: center;">Accept Invitation</a>
									</td>
								</tr>
							</table>
							<!-- Alternative Link -->
							<p style="margin: 0 0 8px; color: #666666; font-size: 14px; line-height: 1.6;">Or copy and paste this link into your browser:</p>
							<p style="margin: 0 0 24px; word-break: break-all; color: #27ae60; font-size: 14px; line-height: 1.6;">
								<a href="%s" style="color: #27ae60; text-decoration: underline;">%s</a>
							</p>
							<!-- Expiration Notice -->
							<div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 12px 16px; margin: 24px 0; border-radius: 4px;">
								<p style="margin: 0; color: #856404; font-size: 14px; line-height: 1.6;">
									<strong>Note:</strong> This invitation will expire in 7 days. Please accept it soon.
								</p>
							</div>
						</td>
					</tr>
					<!-- Footer -->
					<tr>
						<td style="padding: 24px 40px; background-color: #f8f9fa; border-top: 1px solid #e9ecef; border-radius: 0 0 8px 8px; text-align: center;">
							<p style="margin: 0; color: #6c757d; font-size: 12px; line-height: 1.6;">
								This is an automated message from %s. Please do not reply to this email.
							</p>
							<p style="margin: 8px 0 0; color: #6c757d; font-size: 12px; line-height: 1.6;">
								If you did not expect this invitation, you can safely ignore this email.
							</p>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>
`, data.CompanyName, data.InviterName, data.CompanyName, invitationURL, invitationURL, invitationURL, data.CompanyName)

	plainBody = fmt.Sprintf(`
You're Invited!

Hello,

%s has invited you to join %s.

Click the following link to accept the invitation and set up your account:
%s

Note: This invitation will expire in 7 days. Please accept it soon.

---
This is an automated message from %s. Please do not reply to this email.
If you did not expect this invitation, you can safely ignore this email.
`, data.InviterName, data.CompanyName, invitationURL, data.CompanyName)

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

// GenerateCredentialsEmail generates HTML and plain text versions of credentials email
func GenerateCredentialsEmail(data EmailTemplateData) (htmlBody, plainBody string) {
	password := ""
	inviterName := "A team member"
	loginURL := "http://localhost:3000/login"

	if data.CustomData != nil {
		if pwd, ok := data.CustomData["password"].(string); ok {
			password = pwd
		}
		if inv, ok := data.CustomData["inviterName"].(string); ok && inv != "" {
			inviterName = inv
		}
		if url, ok := data.CustomData["loginURL"].(string); ok && url != "" {
			loginURL = url
		}
	}

	htmlBody = fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>Welcome to %s</title>
	<!--[if mso]>
	<style type="text/css">
		table {border-collapse:collapse;border-spacing:0;margin:0;}
		div, td {padding:0;}
		div {margin:0 !important;}
	</style>
	<![endif]-->
</head>
<body style="margin: 0; padding: 0; background-color: #f4f4f4; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;">
	<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin: 0; padding: 0; background-color: #f4f4f4;">
		<tr>
			<td align="center" style="padding: 40px 20px;">
				<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="max-width: 600px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
					<!-- Header -->
					<tr>
						<td style="padding: 40px 40px 20px; text-align: center; background-color: #ffffff; border-radius: 8px 8px 0 0;">
							<h1 style="margin: 0; color: #2c3e50; font-size: 28px; font-weight: 600;">Welcome to %s!</h1>
						</td>
					</tr>
					<!-- Content -->
					<tr>
						<td style="padding: 20px 40px;">
							<p style="margin: 0 0 16px; color: #333333; font-size: 16px; line-height: 1.6;">Hello,</p>
							<p style="margin: 0 0 16px; color: #333333; font-size: 16px; line-height: 1.6;">
								<strong>%s</strong> has added you to <strong>%s</strong>. Your account has been created with the following credentials:
							</p>
							<!-- Credentials Box -->
							<div style="background-color: #f8f9fa; border-left: 4px solid #3498db; padding: 16px; margin: 24px 0; border-radius: 4px;">
								<p style="margin: 0 0 12px; color: #333333; font-size: 14px; line-height: 1.6;">
									<strong>Email:</strong> <span style="color: #2c3e50;">%s</span>
								</p>
								<p style="margin: 0; color: #333333; font-size: 14px; line-height: 1.6;">
									<strong>Password:</strong> <code style="background-color: #e9ecef; padding: 4px 8px; border-radius: 3px; font-family: 'Courier New', monospace; font-size: 13px; color: #2c3e50;">%s</code>
								</p>
							</div>
							<!-- Security Notice -->
							<div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 12px 16px; margin: 24px 0; border-radius: 4px;">
								<p style="margin: 0; color: #856404; font-size: 14px; line-height: 1.6;">
									<strong>Important:</strong> Please change your password after your first login for security.
								</p>
							</div>
							<!-- CTA Button -->
							<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%">
								<tr>
									<td align="center" style="padding: 0 0 24px;">
										<a href="%s" style="display: inline-block; background-color: #3498db; color: #ffffff; text-decoration: none; padding: 14px 32px; border-radius: 6px; font-size: 16px; font-weight: 600; text-align: center;">Login to Your Account</a>
									</td>
								</tr>
							</table>
							<!-- Alternative Link -->
							<p style="margin: 0 0 8px; color: #666666; font-size: 14px; line-height: 1.6;">Or copy and paste this link into your browser:</p>
							<p style="margin: 0 0 24px; word-break: break-all; color: #3498db; font-size: 14px; line-height: 1.6;">
								<a href="%s" style="color: #3498db; text-decoration: underline;">%s</a>
							</p>
						</td>
					</tr>
					<!-- Footer -->
					<tr>
						<td style="padding: 24px 40px; background-color: #f8f9fa; border-top: 1px solid #e9ecef; border-radius: 0 0 8px 8px; text-align: center;">
							<p style="margin: 0; color: #6c757d; font-size: 12px; line-height: 1.6;">
								This is an automated message from %s. Please do not reply to this email.
							</p>
							<p style="margin: 8px 0 0; color: #6c757d; font-size: 12px; line-height: 1.6;">
								Keep your credentials secure and do not share them with anyone.
							</p>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>
`, data.CompanyName, data.CompanyName, inviterName, data.CompanyName, data.UserName, password, loginURL, loginURL, loginURL, data.CompanyName)

	plainBody = fmt.Sprintf(`
Welcome to %s!

Hello %s,

%s has added you to %s. Your account has been created with the following credentials:

Email: %s
Password: %s

Important: Please change your password after your first login for security.

Login to your account:
%s

--- 
This is an automated message, please do not reply.
`, data.CompanyName, data.UserName, inviterName, data.CompanyName, data.UserName, password, loginURL)

	return htmlBody, plainBody
}

// GenerateWelcomeEmail generates HTML and plain text versions of welcome email for existing users
func GenerateWelcomeEmail(data EmailTemplateData) (htmlBody, plainBody string) {
	loginURL := ""
	if data.CustomData != nil {
		if url, ok := data.CustomData["loginURL"].(string); ok {
			loginURL = url
		}
	}
	if loginURL == "" {
		loginURL = "http://localhost:3000/login"
	}

	htmlBody = fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>Welcome to %s</title>
	<!--[if mso]>
	<style type="text/css">
		table {border-collapse:collapse;border-spacing:0;margin:0;}
		div, td {padding:0;}
		div {margin:0 !important;}
	</style>
	<![endif]-->
</head>
<body style="margin: 0; padding: 0; background-color: #f4f4f4; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;">
	<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin: 0; padding: 0; background-color: #f4f4f4;">
		<tr>
			<td align="center" style="padding: 40px 20px;">
				<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="max-width: 600px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
					<!-- Header -->
					<tr>
						<td style="padding: 40px 40px 20px; text-align: center; background-color: #ffffff; border-radius: 8px 8px 0 0;">
							<h1 style="margin: 0; color: #2c3e50; font-size: 28px; font-weight: 600;">Welcome to %s!</h1>
						</td>
					</tr>
					<!-- Content -->
					<tr>
						<td style="padding: 20px 40px;">
							<p style="margin: 0 0 16px; color: #333333; font-size: 16px; line-height: 1.6;">Hello %s,</p>
							<p style="margin: 0 0 16px; color: #333333; font-size: 16px; line-height: 1.6;">
								<strong>%s</strong> has added you to <strong>%s</strong>.
							</p>
							<!-- Role Info -->
							<div style="background-color: #e8f5e9; border-left: 4px solid #4caf50; padding: 16px; margin: 24px 0; border-radius: 4px;">
								<p style="margin: 0; color: #2e7d32; font-size: 14px; line-height: 1.6;">
									<strong>Your Role:</strong> <span style="text-transform: capitalize;">%s</span>
								</p>
							</div>
							<p style="margin: 0 0 24px; color: #333333; font-size: 16px; line-height: 1.6;">
								You can now access the company dashboard and start working with your team.
							</p>
							<!-- CTA Button -->
							<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%">
								<tr>
									<td align="center" style="padding: 0 0 24px;">
										<a href="%s" style="display: inline-block; background-color: #3498db; color: #ffffff; text-decoration: none; padding: 14px 32px; border-radius: 6px; font-size: 16px; font-weight: 600; text-align: center;">Access Dashboard</a>
									</td>
								</tr>
							</table>
							<!-- Alternative Link -->
							<p style="margin: 0 0 8px; color: #666666; font-size: 14px; line-height: 1.6;">Or copy and paste this link into your browser:</p>
							<p style="margin: 0 0 24px; word-break: break-all; color: #3498db; font-size: 14px; line-height: 1.6;">
								<a href="%s" style="color: #3498db; text-decoration: underline;">%s</a>
							</p>
						</td>
					</tr>
					<!-- Footer -->
					<tr>
						<td style="padding: 24px 40px; background-color: #f8f9fa; border-top: 1px solid #e9ecef; border-radius: 0 0 8px 8px; text-align: center;">
							<p style="margin: 0; color: #6c757d; font-size: 12px; line-height: 1.6;">
								This is an automated message from %s. Please do not reply to this email.
							</p>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>
`, data.CompanyName, data.CompanyName, data.UserName, data.InviterName, data.CompanyName, data.Role, loginURL, loginURL, loginURL, data.CompanyName)

	plainBody = fmt.Sprintf(`
Welcome to %s!

Hello %s,

%s has added you to %s.

Your Role: %s

You can now access the company dashboard and start working with your team.

Access Dashboard:
%s

---
This is an automated message from %s. Please do not reply to this email.
`, data.CompanyName, data.UserName, data.InviterName, data.CompanyName, data.Role, loginURL, data.CompanyName)

	return htmlBody, plainBody
}

// GenerateNewUserSetupEmail generates HTML and plain text versions of new user setup email with OTP
func GenerateNewUserSetupEmail(data EmailTemplateData) (htmlBody, plainBody string) {
	password := ""
	inviterName := "A team member"
	setupURL := data.SetupURL
	if setupURL == "" && data.OTPCode != "" {
		setupURL = fmt.Sprintf("http://localhost:3000/setup-account?otp=%s&email=%s", data.OTPCode, data.UserName)
	}

	if data.CustomData != nil {
		if pwd, ok := data.CustomData["password"].(string); ok {
			password = pwd
		}
		if inv, ok := data.CustomData["inviterName"].(string); ok && inv != "" {
			inviterName = inv
		}
		if url, ok := data.CustomData["setupURL"].(string); ok && url != "" {
			setupURL = url
		}
	}

	htmlBody = fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>Welcome to %s - Complete Your Account Setup</title>
	<!--[if mso]>
	<style type="text/css">
		table {border-collapse:collapse;border-spacing:0;margin:0;}
		div, td {padding:0;}
		div {margin:0 !important;}
	</style>
	<![endif]-->
</head>
<body style="margin: 0; padding: 0; background-color: #f4f4f4; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;">
	<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin: 0; padding: 0; background-color: #f4f4f4;">
		<tr>
			<td align="center" style="padding: 40px 20px;">
				<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="600" style="max-width: 600px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
					<!-- Header -->
					<tr>
						<td style="padding: 40px 40px 20px; text-align: center; background-color: #ffffff; border-radius: 8px 8px 0 0;">
							<h1 style="margin: 0; color: #2c3e50; font-size: 28px; font-weight: 600;">Welcome to %s!</h1>
						</td>
					</tr>
					<!-- Content -->
					<tr>
						<td style="padding: 20px 40px;">
							<p style="margin: 0 0 16px; color: #333333; font-size: 16px; line-height: 1.6;">Hello,</p>
							<p style="margin: 0 0 16px; color: #333333; font-size: 16px; line-height: 1.6;">
								<strong>%s</strong> has added you to <strong>%s</strong>. Your account has been created with the following credentials:
							</p>
							<!-- Credentials Box -->
							<div style="background-color: #f8f9fa; border-left: 4px solid #3498db; padding: 16px; margin: 24px 0; border-radius: 4px;">
								<p style="margin: 0 0 12px; color: #333333; font-size: 14px; line-height: 1.6;">
									<strong>Email:</strong> <span style="color: #2c3e50;">%s</span>
								</p>
								<p style="margin: 0; color: #333333; font-size: 14px; line-height: 1.6;">
									<strong>Password:</strong> <code style="background-color: #e9ecef; padding: 4px 8px; border-radius: 3px; font-family: 'Courier New', monospace; font-size: 13px; color: #2c3e50;">%s</code>
								</p>
							</div>
							<!-- Security Notice -->
							<div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 12px 16px; margin: 24px 0; border-radius: 4px;">
								<p style="margin: 0; color: #856404; font-size: 14px; line-height: 1.6;">
									<strong>Important:</strong> Please complete your account setup by clicking the button below. You'll need to:
								</p>
								<ul style="margin: 8px 0 0 0; padding-left: 20px; color: #856404; font-size: 14px; line-height: 1.6;">
									<li>Change your password</li>
									<li>Add your personal information (name, phone number)</li>
								</ul>
							</div>
							<!-- OTP Code Display -->
							<div style="background-color: #e3f2fd; border-left: 4px solid #2196f3; padding: 16px; margin: 24px 0; border-radius: 4px; text-align: center;">
								<p style="margin: 0 0 8px; color: #1565c0; font-size: 12px; line-height: 1.6; text-transform: uppercase; letter-spacing: 1px;">
									Your Setup Code
								</p>
								<p style="margin: 0; color: #1565c0; font-size: 32px; font-weight: 600; letter-spacing: 4px; font-family: 'Courier New', monospace;">
									%s
								</p>
							</div>
							<!-- CTA Button -->
							<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%">
								<tr>
									<td align="center" style="padding: 0 0 24px;">
										<a href="%s" style="display: inline-block; background-color: #27ae60; color: #ffffff; text-decoration: none; padding: 14px 32px; border-radius: 6px; font-size: 16px; font-weight: 600; text-align: center;">Complete Account Setup</a>
									</td>
								</tr>
							</table>
							<!-- Alternative Link -->
							<p style="margin: 0 0 8px; color: #666666; font-size: 14px; line-height: 1.6;">Or copy and paste this link into your browser:</p>
							<p style="margin: 0 0 24px; word-break: break-all; color: #27ae60; font-size: 14px; line-height: 1.6;">
								<a href="%s" style="color: #27ae60; text-decoration: underline;">%s</a>
							</p>
							<!-- Expiration Notice -->
							<div style="background-color: #ffebee; border-left: 4px solid #f44336; padding: 12px 16px; margin: 24px 0; border-radius: 4px;">
								<p style="margin: 0; color: #c62828; font-size: 14px; line-height: 1.6;">
									<strong>Note:</strong> This setup link will expire in 24 hours. Please complete your setup soon.
								</p>
							</div>
						</td>
					</tr>
					<!-- Footer -->
					<tr>
						<td style="padding: 24px 40px; background-color: #f8f9fa; border-top: 1px solid #e9ecef; border-radius: 0 0 8px 8px; text-align: center;">
							<p style="margin: 0; color: #6c757d; font-size: 12px; line-height: 1.6;">
								This is an automated message from %s. Please do not reply to this email.
							</p>
							<p style="margin: 8px 0 0; color: #6c757d; font-size: 12px; line-height: 1.6;">
								Keep your credentials secure and do not share them with anyone.
							</p>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>
`, data.CompanyName, data.CompanyName, inviterName, data.CompanyName, data.UserName, password, data.OTPCode, setupURL, setupURL, setupURL, data.CompanyName)

	plainBody = fmt.Sprintf(`
Welcome to %s - Complete Your Account Setup

Hello,

%s has added you to %s. Your account has been created with the following credentials:

Email: %s
Password: %s

Your Setup Code: %s

Important: Please complete your account setup by visiting:
%s

You'll need to:
- Change your password
- Add your personal information (name, phone number)

Note: This setup link will expire in 24 hours. Please complete your setup soon.

---
This is an automated message from %s. Please do not reply to this email.
Keep your credentials secure and do not share them with anyone.
`, data.CompanyName, inviterName, data.CompanyName, data.UserName, password, data.OTPCode, setupURL, data.CompanyName)

	return htmlBody, plainBody
}
