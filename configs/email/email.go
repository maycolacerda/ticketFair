// configs/email.go
package configs

import (
	"log/slog"
	"os"
	"strconv"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

var Email *EmailConfig

func InitEmail() {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		slog.Error("Invalid SMTP_PORT", "error", err.Error())
		port = 587 // ← default TLS port
	}

	Email = &EmailConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
	}

	missing := []string{}
	if Email.Host == "" {
		missing = append(missing, "SMTP_HOST")
	}
	if Email.Username == "" {
		missing = append(missing, "SMTP_USERNAME")
	}
	if Email.Password == "" {
		missing = append(missing, "SMTP_PASSWORD")
	}
	if Email.From == "" {
		missing = append(missing, "SMTP_FROM")
	}

	if len(missing) > 0 {
		slog.Warn("Missing email configuration — email sending disabled", "vars", missing)
		Email = nil // ← nil signals email is disabled
	}
}
