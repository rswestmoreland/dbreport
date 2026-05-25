package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rswestmoreland/dbreport/internal/config"
)

func TestHelpCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	r := Runner{Stdout: &stdout, Stderr: &stderr}

	code := r.Run([]string{"help"})
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("expected help output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestSampleConfigCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	r := Runner{Stdout: &stdout, Stderr: &stderr}

	code := r.Run([]string{"sample-config"})
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "title:") || !strings.Contains(stdout.String(), "queries:") {
		t.Fatalf("expected sample config output, got %q", stdout.String())
	}
}

func TestCheckCommandInvalidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.yml")
	if err := os.WriteFile(path, []byte("title: ''\n"), 0600); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	r := Runner{Stdout: &stdout, Stderr: &stderr}

	code := r.Run([]string{"check", "--config", path})
	if code != ExitConfigError {
		t.Fatalf("expected config error, got %d", code)
	}
	if !strings.Contains(stderr.String(), "configuration validation failed") {
		t.Fatalf("expected validation error, got %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}
}

func TestCheckCommandMissingDatabaseEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.yml")
	if err := os.WriteFile(path, []byte(config.SampleYAML), 0600); err != nil {
		t.Fatalf("write sample config: %v", err)
	}

	t.Setenv("DBREPORT_DB_USER", "")
	t.Setenv("DBREPORT_DB_PASSWORD", "")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	r := Runner{Stdout: &stdout, Stderr: &stderr}

	code := r.Run([]string{"check", "--config", path})
	if code != ExitDatabaseError {
		t.Fatalf("expected database error, got %d", code)
	}
	if !strings.Contains(stderr.String(), "database username environment variable") {
		t.Fatalf("expected missing env error, got %q", stderr.String())
	}
}

func TestRunCommandEmailValidatesConfigBeforeDatabase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.yml")
	yaml := strings.Replace(config.SampleYAML, `smtp_host: "smtp.example.com"`, `smtp_host: ""`, 1)
	if err := os.WriteFile(path, []byte(yaml), 0600); err != nil {
		t.Fatalf("write sample config: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	r := Runner{Stdout: &stdout, Stderr: &stderr}

	code := r.Run([]string{"run", "--config", path, "--email"})
	if code != ExitConfigError {
		t.Fatalf("expected config error, got %d", code)
	}
	if !strings.Contains(stderr.String(), "email.smtp_host is required") {
		t.Fatalf("expected email config validation error, got %q", stderr.String())
	}
}

func TestParseCommonOptionsRejectsQuietVerboseTogether(t *testing.T) {
	_, err := parseCommonOptions([]string{"--quiet", "--verbose"})
	if err == nil {
		t.Fatalf("expected quiet/verbose conflict")
	}
	if !strings.Contains(err.Error(), "--quiet and --verbose") {
		t.Fatalf("expected quiet/verbose error, got %q", err.Error())
	}
}

func TestParseCommonOptionsRejectsEmailNoEmailTogether(t *testing.T) {
	_, err := parseCommonOptions([]string{"--email", "--no-email"})
	if err == nil {
		t.Fatalf("expected email/no-email conflict")
	}
	if !strings.Contains(err.Error(), "--email and --no-email") {
		t.Fatalf("expected email/no-email error, got %q", err.Error())
	}
}

func TestCheckCommandRejectsOutputOption(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	r := Runner{Stdout: &stdout, Stderr: &stderr}

	code := r.Run([]string{"check", "--output", "ignored.html"})
	if code != ExitConfigError {
		t.Fatalf("expected config error, got %d", code)
	}
	if !strings.Contains(stderr.String(), "--output is only valid with run") {
		t.Fatalf("expected output option error, got %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}
}

func TestRunHelpCommandSpecificOutput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	r := Runner{Stdout: &stdout, Stderr: &stderr}

	code := r.Run([]string{"run", "--help"})
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d", code)
	}
	if !strings.Contains(stdout.String(), "--no-email") || !strings.Contains(stdout.String(), "--verbose") {
		t.Fatalf("expected run help details, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestVersionCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	r := Runner{Stdout: &stdout, Stderr: &stderr}

	code := r.Run([]string{"version"})
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d", code)
	}
	if !strings.Contains(stdout.String(), "dbreport") || !strings.Contains(stdout.String(), "Build date:") {
		t.Fatalf("expected version output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestAboutCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	r := Runner{Stdout: &stdout, Stderr: &stderr}

	code := r.Run([]string{"about"})
	if code != ExitSuccess {
		t.Fatalf("expected success, got %d", code)
	}
	for _, want := range []string{"Richard S. Westmoreland", "dev@rswestmore.land", "MIT", "Copyright"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected about output to contain %q, got %q", want, stdout.String())
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
}

func TestUnknownCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	r := Runner{Stdout: &stdout, Stderr: &stderr}

	code := r.Run([]string{"bogus"})
	if code != ExitGeneralError {
		t.Fatalf("expected general error, got %d", code)
	}
	if !strings.Contains(stderr.String(), "unknown command: bogus") {
		t.Fatalf("expected unknown command error, got %q", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}
}
