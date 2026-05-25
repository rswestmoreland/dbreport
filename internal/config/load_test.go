package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFileReadsYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.yml")
	if err := os.WriteFile(path, []byte(SampleYAML), 0600); err != nil {
		t.Fatalf("write sample config: %v", err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}
	if cfg.Title != "Daily MariaDB Summary" {
		t.Fatalf("unexpected title: %q", cfg.Title)
	}
	if len(cfg.Queries) != 4 {
		t.Fatalf("expected 4 sample queries, got %d", len(cfg.Queries))
	}
}

func TestLoadFileRejectsInvalidYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.yml")
	if err := os.WriteFile(path, []byte("title: [unterminated\n"), 0600); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	_, err := LoadFile(path)
	if err == nil {
		t.Fatal("expected parse error")
	}
	if !strings.Contains(err.Error(), "parse config") {
		t.Fatalf("expected parse config error, got %v", err)
	}
}

func TestValidateForEmailRequestRequiresEmailConfig(t *testing.T) {
	cfg := validConfigForTest()
	cfg.Email = EmailConfig{}

	err := cfg.ValidateForEmailRequest()
	if err == nil {
		t.Fatal("expected email validation error")
	}
	if !strings.Contains(err.Error(), "email.smtp_host") {
		t.Fatalf("expected email smtp_host error, got %v", err)
	}
}

func TestValidateRejectsEmailHeaderBreaks(t *testing.T) {
	cfg := validConfigForTest()
	cfg.Email = EmailConfig{
		Enabled:      true,
		SMTPHost:     "smtp.example.com",
		SMTPPort:     587,
		From:         "reports@example.com",
		To:           []string{"ops@example.com"},
		Subject:      "Daily Report\nBCC: attacker@example.com",
		SendHTMLBody: true,
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected email header validation error")
	}
	if !strings.Contains(err.Error(), "email.subject must not contain newlines") {
		t.Fatalf("expected subject newline error, got %v", err)
	}
}
