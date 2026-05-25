package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/rswestmoreland/dbreport/internal/config"
)

type SendRequest struct {
	Config      config.EmailConfig
	HTML        []byte
	ReportPath  string
	GeneratedAt time.Time
}

func Send(req SendRequest) error {
	cfg := req.Config
	if strings.TrimSpace(cfg.SMTPHost) == "" {
		return fmt.Errorf("email.smtp_host is required")
	}
	if cfg.SMTPPort <= 0 || cfg.SMTPPort > 65535 {
		return fmt.Errorf("email.smtp_port must be between 1 and 65535")
	}
	if req.GeneratedAt.IsZero() {
		req.GeneratedAt = time.Now()
	}

	message, err := BuildMessage(cfg, req.HTML, req.ReportPath, req.GeneratedAt)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(cfg.SMTPHost, fmt.Sprint(cfg.SMTPPort))
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("connect to SMTP server %s: %w", addr, err)
	}
	defer client.Close()

	if cfg.StartTLS {
		tlsConfig := &tls.Config{ServerName: cfg.SMTPHost, MinVersion: tls.VersionTLS12}
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return fmt.Errorf("SMTP server %s does not advertise STARTTLS", addr)
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("start SMTP TLS: %w", err)
		}
	}

	username, password, useAuth, err := smtpCredentials(cfg)
	if err != nil {
		return err
	}
	if useAuth {
		auth := smtp.PlainAuth("", username, password, cfg.SMTPHost)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("authenticate to SMTP server: %w", err)
		}
	}

	from, err := envelopeAddress(cfg.From)
	if err != nil {
		return fmt.Errorf("parse SMTP sender: %w", err)
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("set SMTP sender: %w", err)
	}
	for _, recipient := range cfg.To {
		to, err := envelopeAddress(recipient)
		if err != nil {
			return fmt.Errorf("parse SMTP recipient %q: %w", recipient, err)
		}
		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("set SMTP recipient %q: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("start SMTP DATA: %w", err)
	}
	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return fmt.Errorf("write SMTP message: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("finish SMTP message: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("quit SMTP session: %w", err)
	}
	return nil
}

func smtpCredentials(cfg config.EmailConfig) (string, string, bool, error) {
	userEnv := strings.TrimSpace(cfg.UsernameEnv)
	passEnv := strings.TrimSpace(cfg.PasswordEnv)
	if userEnv == "" && passEnv == "" {
		return "", "", false, nil
	}
	if userEnv == "" || passEnv == "" {
		return "", "", false, fmt.Errorf("email.username_env and email.password_env must be set together")
	}

	username := os.Getenv(userEnv)
	if username == "" {
		return "", "", false, fmt.Errorf("SMTP username environment variable %s is not set", userEnv)
	}
	password := os.Getenv(passEnv)
	if password == "" {
		return "", "", false, fmt.Errorf("SMTP password environment variable %s is not set", passEnv)
	}
	return username, password, true, nil
}
