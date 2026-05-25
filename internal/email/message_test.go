package email

import (
	"strings"
	"testing"
	"time"

	"github.com/rswestmoreland/dbreport/internal/config"
)

func TestBuildMessageHTMLBody(t *testing.T) {
	cfg := config.EmailConfig{
		From:         "reports@example.com",
		To:           []string{"ops@example.com"},
		Subject:      "Daily Report",
		SendHTMLBody: true,
	}

	msg, err := BuildMessage(cfg, []byte("<html><body>ok</body></html>"), "report.html", fixedTime())
	if err != nil {
		t.Fatalf("BuildMessage failed: %v", err)
	}
	text := string(msg)
	for _, want := range []string{
		"From: reports@example.com",
		"To: ops@example.com",
		"Subject: Daily Report",
		"Content-Type: multipart/alternative",
		"Content-Type: text/plain; charset=utf-8",
		"Content-Type: text/html; charset=utf-8",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected message to contain %q, got:\n%s", want, text)
		}
	}
}

func TestBuildMessageAttachment(t *testing.T) {
	cfg := config.EmailConfig{
		From:       "reports@example.com",
		To:         []string{"ops@example.com"},
		Subject:    "Daily Report",
		AttachHTML: true,
	}

	msg, err := BuildMessage(cfg, []byte("<html><body>ok</body></html>"), "reports/daily.html", fixedTime())
	if err != nil {
		t.Fatalf("BuildMessage failed: %v", err)
	}
	text := string(msg)
	for _, want := range []string{
		"Content-Type: multipart/mixed",
		"Content-Disposition: attachment",
		"daily.html",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected message to contain %q, got:\n%s", want, text)
		}
	}
}

func TestBuildMessageRejectsHeaderInjection(t *testing.T) {
	cfg := config.EmailConfig{
		From:         "reports@example.com",
		To:           []string{"ops@example.com"},
		Subject:      "Daily\nBCC: attacker@example.com",
		SendHTMLBody: true,
	}

	_, err := BuildMessage(cfg, []byte("<html></html>"), "report.html", fixedTime())
	if err == nil {
		t.Fatal("expected header injection error")
	}
}

func fixedTime() time.Time {
	return time.Date(2026, 5, 25, 12, 0, 0, 0, time.UTC)
}
