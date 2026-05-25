# Changelog

## 0.1.0-dev

### Added

- Project skeleton.
- MIT license.
- Author, license, copyright, and build metadata.
- Usage-focused help output.
- `version` command.
- `about` command.
- YAML configuration with `github.com/goccy/go-yaml`.
- Sample configuration output.
- Config loading and validation.
- CLI option parsing for `--config`, `--output`, `--email`, `--no-email`,
  `--quiet`, and `--verbose`.
- MariaDB/MySQL driver integration through Go `database/sql`.
- Database connection handling using environment-based credentials.
- Query execution with per-query timeout.
- Row cap enforcement.
- Generic query result model.
- Self-contained HTML report rendering.
- Metric, table, bar chart, and line chart sections.
- Inline CSS and inline SVG charts.
- Optional SMTP email support using Go standard-library networking packages.
- HTML body email mode.
- Optional HTML attachment mode.
- Plain-text email fallback.
- STARTTLS support.
- Optional SMTP authentication through environment variables.
- Command, configuration, report section, email, security, build, and release docs.
- Example YAML configurations.
- Release packaging scripts and Makefile targets.
