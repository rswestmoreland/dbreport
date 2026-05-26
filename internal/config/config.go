package config

import (
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

const DefaultConfigPath = "report.yml"

type Config struct {
	Title     string          `yaml:"title"`
	Output    OutputConfig    `yaml:"output"`
	Database  DatabaseConfig  `yaml:"database"`
	Limits    LimitsConfig    `yaml:"limits"`
	Rendering RenderingConfig `yaml:"rendering"`
	Safety    SafetyConfig    `yaml:"safety"`
	Email     EmailConfig     `yaml:"email"`
	Queries   []QueryConfig   `yaml:"queries"`
}

type OutputConfig struct {
	HTML string `yaml:"html"`
}

type DatabaseConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	Name           string `yaml:"name"`
	UserEnv        string `yaml:"user_env"`
	PasswordEnv    string `yaml:"password_env"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	TLS            bool   `yaml:"tls"`
}

type LimitsConfig struct {
	MaxRowsPerQuery int `yaml:"max_rows_per_query"`
	MaxCellLength   int `yaml:"max_cell_length"`
	MaxReportBytes  int `yaml:"max_report_bytes"`
}

type RenderingConfig struct {
	NullValue string `yaml:"null_value"`
}

type SafetyConfig struct {
	AllowedDatabases []string `yaml:"allowed_databases"`
	AllowedTables    []string `yaml:"allowed_tables"`
	BlockedFunctions []string `yaml:"blocked_functions"`
	BlockedColumns   []string `yaml:"blocked_columns"`
	BlockedPatterns  []string `yaml:"blocked_patterns"`
}

type EmailConfig struct {
	Enabled      bool     `yaml:"enabled"`
	SMTPHost     string   `yaml:"smtp_host"`
	SMTPPort     int      `yaml:"smtp_port"`
	StartTLS     bool     `yaml:"starttls"`
	UsernameEnv  string   `yaml:"username_env"`
	PasswordEnv  string   `yaml:"password_env"`
	From         string   `yaml:"from"`
	To           []string `yaml:"to"`
	Subject      string   `yaml:"subject"`
	SendHTMLBody bool     `yaml:"send_html_body"`
	AttachHTML   bool     `yaml:"attach_html"`
}

type QueryConfig struct {
	NullValue    string `yaml:"null_value"`
	ID           string `yaml:"id"`
	Title        string `yaml:"title"`
	Type         string `yaml:"type"`
	LabelColumn  string `yaml:"label_column"`
	SeriesColumn string `yaml:"series_column"`
	ValueColumn  string `yaml:"value_column"`
	ShowTable    *bool  `yaml:"show_table"`
	SQL          string `yaml:"sql"`
}

func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	return c.validate(false)
}

func (c *Config) ValidateForEmailRequest() error {
	return c.validate(true)
}

func (c *Config) validate(emailRequested bool) error {
	var problems []string

	if strings.TrimSpace(c.Title) == "" {
		problems = append(problems, "title is required")
	}
	if strings.TrimSpace(c.Output.HTML) == "" {
		problems = append(problems, "output.html is required")
	}

	validateDatabase(c.Database, &problems)
	validateLimits(c.Limits, &problems)
	validateQueries(c.Queries, &problems)
	if c.Email.Enabled || emailRequested {
		validateEmail(c.Email, &problems)
	}

	if len(problems) > 0 {
		return ValidationError{Problems: problems}
	}
	return nil
}

func validateDatabase(db DatabaseConfig, problems *[]string) {
	if strings.TrimSpace(db.Host) == "" {
		*problems = append(*problems, "database.host is required")
	}
	if db.Port <= 0 || db.Port > 65535 {
		*problems = append(*problems, "database.port must be between 1 and 65535")
	}
	if strings.TrimSpace(db.Name) == "" {
		*problems = append(*problems, "database.name is required")
	}
	if strings.TrimSpace(db.UserEnv) == "" {
		*problems = append(*problems, "database.user_env is required")
	}
	if strings.TrimSpace(db.PasswordEnv) == "" {
		*problems = append(*problems, "database.password_env is required")
	}
	if db.TimeoutSeconds <= 0 {
		*problems = append(*problems, "database.timeout_seconds must be greater than zero")
	}
}

func validateLimits(limits LimitsConfig, problems *[]string) {
	if limits.MaxRowsPerQuery <= 0 {
		*problems = append(*problems, "limits.max_rows_per_query must be greater than zero")
	}
	if limits.MaxCellLength < 0 {
		*problems = append(*problems, "limits.max_cell_length must be zero or greater")
	}
	if limits.MaxReportBytes < 0 {
		*problems = append(*problems, "limits.max_report_bytes must be zero or greater")
	}
}

func validateQueries(queries []QueryConfig, problems *[]string) {
	if len(queries) == 0 {
		*problems = append(*problems, "queries must contain at least one query")
		return
	}

	seen := make(map[string]struct{}, len(queries))
	for i, q := range queries {
		prefix := fmt.Sprintf("queries[%d]", i)
		id := strings.TrimSpace(q.ID)
		qtype := strings.TrimSpace(q.Type)

		if id == "" {
			*problems = append(*problems, prefix+".id is required")
		} else if _, ok := seen[id]; ok {
			*problems = append(*problems, prefix+".id is duplicated: "+id)
		} else {
			seen[id] = struct{}{}
		}

		if strings.TrimSpace(q.Title) == "" {
			*problems = append(*problems, prefix+".title is required")
		}
		if qtype == "" {
			*problems = append(*problems, prefix+".type is required")
		} else if !isKnownQueryType(qtype) {
			*problems = append(*problems, prefix+".type is unsupported: "+qtype)
		}
		if strings.TrimSpace(q.SQL) == "" {
			*problems = append(*problems, prefix+".sql is required")
		}
		if qtype == "bar" || qtype == "line" || qtype == "pie" {
			if strings.TrimSpace(q.LabelColumn) == "" {
				*problems = append(*problems, prefix+".label_column is required for "+qtype+" sections")
			}
			if strings.TrimSpace(q.ValueColumn) == "" {
				*problems = append(*problems, prefix+".value_column is required for "+qtype+" sections")
			}
		}
	}
}

func validateEmail(email EmailConfig, problems *[]string) {
	if strings.TrimSpace(email.SMTPHost) == "" {
		*problems = append(*problems, "email.smtp_host is required when email is enabled or requested")
	}
	if email.SMTPPort <= 0 || email.SMTPPort > 65535 {
		*problems = append(*problems, "email.smtp_port must be between 1 and 65535 when email is enabled or requested")
	}
	from := strings.TrimSpace(email.From)
	if from == "" {
		*problems = append(*problems, "email.from is required when email is enabled or requested")
	} else if containsHeaderBreak(from) {
		*problems = append(*problems, "email.from must not contain newlines")
	} else if _, err := mail.ParseAddress(from); err != nil {
		*problems = append(*problems, "email.from is invalid: "+err.Error())
	}

	if len(email.To) == 0 {
		*problems = append(*problems, "email.to must contain at least one recipient when email is enabled or requested")
	}
	for i, recipient := range email.To {
		trimmed := strings.TrimSpace(recipient)
		if trimmed == "" {
			*problems = append(*problems, fmt.Sprintf("email.to[%d] is required", i))
		} else if containsHeaderBreak(trimmed) {
			*problems = append(*problems, fmt.Sprintf("email.to[%d] must not contain newlines", i))
		} else if _, err := mail.ParseAddress(trimmed); err != nil {
			*problems = append(*problems, fmt.Sprintf("email.to[%d] is invalid: %s", i, err.Error()))
		}
	}

	subject := strings.TrimSpace(email.Subject)
	if subject == "" {
		*problems = append(*problems, "email.subject is required when email is enabled or requested")
	} else if containsHeaderBreak(subject) {
		*problems = append(*problems, "email.subject must not contain newlines")
	}
	if (strings.TrimSpace(email.UsernameEnv) == "") != (strings.TrimSpace(email.PasswordEnv) == "") {
		*problems = append(*problems, "email.username_env and email.password_env must be set together")
	}
	if !email.SendHTMLBody && !email.AttachHTML {
		*problems = append(*problems, "email.send_html_body or email.attach_html must be true when email is enabled or requested")
	}
}

func containsHeaderBreak(value string) bool {
	return strings.ContainsAny(value, "\r\n")
}

func isKnownQueryType(qtype string) bool {
	switch qtype {
	case "metric", "table", "bar", "line", "pie":
		return true
	default:
		return false
	}
}

type ValidationError struct {
	Problems []string
}

func (e ValidationError) Error() string {
	var b strings.Builder
	b.WriteString("configuration validation failed")
	for _, problem := range e.Problems {
		b.WriteString("\n  - ")
		b.WriteString(problem)
	}
	return b.String()
}
