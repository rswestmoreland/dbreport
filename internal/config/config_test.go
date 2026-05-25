package config

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestSampleYAMLValidates(t *testing.T) {
	var cfg Config
	if err := yaml.Unmarshal([]byte(SampleYAML), &cfg); err != nil {
		t.Fatalf("sample yaml did not parse: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("sample yaml did not validate: %v", err)
	}
}

func TestValidateRejectsDuplicateQueryID(t *testing.T) {
	cfg := validConfigForTest()
	cfg.Queries = append(cfg.Queries, cfg.Queries[0])

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "duplicated") {
		t.Fatalf("expected duplicate error, got: %v", err)
	}
}

func TestValidateRequiresChartColumns(t *testing.T) {
	cfg := validConfigForTest()
	cfg.Queries = []QueryConfig{{
		ID:    "chart",
		Title: "Chart",
		Type:  "bar",
		SQL:   "SELECT status, COUNT(*) AS count FROM orders GROUP BY status",
	}}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "label_column") || !strings.Contains(err.Error(), "value_column") {
		t.Fatalf("expected chart column errors, got: %v", err)
	}
}

func validConfigForTest() Config {
	return Config{
		Title:  "Test Report",
		Output: OutputConfig{HTML: "report.html"},
		Database: DatabaseConfig{
			Host:           "127.0.0.1",
			Port:           3306,
			Name:           "appdb",
			UserEnv:        "DBREPORT_DB_USER",
			PasswordEnv:    "DBREPORT_DB_PASSWORD",
			TimeoutSeconds: 10,
		},
		Limits: LimitsConfig{MaxRowsPerQuery: 1000},
		Queries: []QueryConfig{{
			ID:    "metric",
			Title: "Metric",
			Type:  "metric",
			SQL:   "SELECT COUNT(*) AS value FROM orders",
		}},
	}
}
