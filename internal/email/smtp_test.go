package email

import (
	"strings"
	"testing"

	"github.com/rswestmoreland/dbreport/internal/config"
)

func TestSMTPCredentialsDisabledWhenEnvNamesEmpty(t *testing.T) {
	username, password, useAuth, err := smtpCredentials(config.EmailConfig{})
	if err != nil {
		t.Fatalf("smtpCredentials failed: %v", err)
	}
	if useAuth {
		t.Fatal("expected auth to be disabled")
	}
	if username != "" || password != "" {
		t.Fatalf("expected empty credentials, got %q/%q", username, password)
	}
}

func TestSMTPCredentialsRequiresBothEnvNames(t *testing.T) {
	_, _, _, err := smtpCredentials(config.EmailConfig{UsernameEnv: "DBREPORT_SMTP_USER"})
	if err == nil {
		t.Fatal("expected paired env validation error")
	}
	if !strings.Contains(err.Error(), "must be set together") {
		t.Fatalf("expected paired env error, got %v", err)
	}
}

func TestSMTPCredentialsReadsEnv(t *testing.T) {
	t.Setenv("DBREPORT_TEST_SMTP_USER", "smtp-user")
	t.Setenv("DBREPORT_TEST_SMTP_PASSWORD", "smtp-password")

	username, password, useAuth, err := smtpCredentials(config.EmailConfig{
		UsernameEnv: "DBREPORT_TEST_SMTP_USER",
		PasswordEnv: "DBREPORT_TEST_SMTP_PASSWORD",
	})
	if err != nil {
		t.Fatalf("smtpCredentials failed: %v", err)
	}
	if !useAuth {
		t.Fatal("expected auth to be enabled")
	}
	if username != "smtp-user" || password != "smtp-password" {
		t.Fatalf("unexpected credentials %q/%q", username, password)
	}
}
