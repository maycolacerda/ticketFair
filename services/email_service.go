// services/email_service.go
package services

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	configs "github.com/maycolacerda/ticketfair/configs/email"
)

func sendEmail(to, subject, body string) error {
	if configs.Email == nil {
		slog.Warn("Email sending skipped — SMTP not configured", "to", to, "subject", subject)
		return nil // ← graceful degradation, don't crash
	}

	auth := smtp.PlainAuth(
		"",
		configs.Email.Username,
		configs.Email.Password,
		configs.Email.Host,
	)

	from := fmt.Sprintf("%s <%s>", configs.Email.FromName, configs.Email.From)

	headers := strings.Join([]string{
		"From: " + from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
	}, "\r\n")

	message := []byte(headers + "\r\n\r\n" + body)

	addr := fmt.Sprintf("%s:%d", configs.Email.Host, configs.Email.Port)

	if err := smtp.SendMail(addr, auth, configs.Email.From, []string{to}, message); err != nil {
		slog.Error("Failed to send email",
			"to", to,
			"subject", subject,
			"error", err.Error(),
		)
		return ErrFailedToCreate
	}

	slog.Info("Email sent", "to", to, "subject", subject)
	return nil
}

func SendVerificationEmail(to, code string) error {
	tmpl := verificationEmailTemplate(code, "15")
	return sendEmail(to, tmpl.Subject, tmpl.Body)
}

func SendWelcomeEmail(to, username string) error {
	tmpl := welcomeEmailTemplate(username)
	return sendEmail(to, tmpl.Subject, tmpl.Body)
}

func SendPurchaseConfirmationEmail(to, username, eventName, ticketID string, amount float64) error {
	tmpl := purchaseConfirmationTemplate(username, eventName, ticketID, amount)
	return sendEmail(to, tmpl.Subject, tmpl.Body)
}

func SendPasswordResetEmail(to, code string) error {
	tmpl := passwordResetEmailTemplate(code, "15")
	return sendEmail(to, tmpl.Subject, tmpl.Body)
}
