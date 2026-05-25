package db

import (
	"strings"
	"testing"

	"github.com/rswestmoreland/dbreport/internal/config"
)

func TestBuildDSN(t *testing.T) {
	cfg := config.DatabaseConfig{
		Host:           "127.0.0.1",
		Port:           3306,
		Name:           "appdb",
		TimeoutSeconds: 10,
		TLS:            true,
	}

	dsn, timeout, err := BuildDSN(cfg, "report_user", "secret")
	if err != nil {
		t.Fatalf("BuildDSN failed: %v", err)
	}
	if timeout.String() != "10s" {
		t.Fatalf("expected timeout 10s, got %s", timeout)
	}
	for _, want := range []string{"report_user:secret@tcp(127.0.0.1:3306)/appdb", "parseTime=true", "timeout=10s", "readTimeout=10s", "writeTimeout=10s", "tls=true"} {
		if !strings.Contains(dsn, want) {
			t.Fatalf("expected DSN to contain %q, got %q", want, dsn)
		}
	}
}

func TestCredentialsFromEnv(t *testing.T) {
	cfg := config.DatabaseConfig{
		UserEnv:     "DBREPORT_TEST_USER",
		PasswordEnv: "DBREPORT_TEST_PASSWORD",
	}
	t.Setenv("DBREPORT_TEST_USER", "user1")
	t.Setenv("DBREPORT_TEST_PASSWORD", "pass1")

	user, password, err := CredentialsFromEnv(cfg)
	if err != nil {
		t.Fatalf("CredentialsFromEnv failed: %v", err)
	}
	if user != "user1" || password != "pass1" {
		t.Fatalf("unexpected credentials %q/%q", user, password)
	}
}

func TestSafeTarget(t *testing.T) {
	cfg := config.DatabaseConfig{Host: "db.example.com", Port: 3306, Name: "reports"}
	if got := SafeTarget(cfg); got != "db.example.com:3306/reports" {
		t.Fatalf("unexpected safe target: %s", got)
	}
}
