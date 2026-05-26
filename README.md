# dbreport

`dbreport` is a small command-line utility for generating clean, self-contained
HTML reports from MariaDB queries.

It is designed for simple operational reporting where a full web application,
BI platform, scheduler, or PDF renderer would be unnecessary.

## Goals

- Single Go binary.
- YAML report configuration.
- MariaDB queries from config files.
- Clean self-contained HTML output.
- Inline CSS and inline SVG charts.
- Optional SMTP email delivery.
- No external runtime dependencies beyond the binary, config file, and
  environment variables.
- Secrets loaded from environment variables, not command-line arguments.
- No arbitrary SQL from the command line.

## Features

- Usage-focused help.
- `version` command.
- `about` command with author, license, copyright, and build metadata.
- YAML config loading and validation.
- MariaDB connection path using Go `database/sql`.
- Query execution with per-query timeout.
- Row cap enforcement.
- Metric, table, bar chart, line chart, and pie chart report sections.
- Self-contained HTML report rendering.
- Optional SMTP email delivery.
- CLI polish flags: `--quiet`, `--verbose`, `--email`, and `--no-email`.

## Commands

```text
dbreport --help
dbreport help
dbreport version
dbreport about
dbreport sample-config
dbreport check --config report.yml
dbreport run --config report.yml
dbreport run --config report.yml --output reports/daily.html
dbreport run --config report.yml --email
dbreport run --config report.yml --no-email
```

## Quick start

Generate a sample configuration:

```sh
dbreport sample-config > report.yml
```

Set database credentials through environment variables:

```sh
export DBREPORT_DB_USER='report_user'
export DBREPORT_DB_PASSWORD='change-me'
```

Validate the configuration and database/query access:

```sh
dbreport check --config report.yml
```

Generate a report:

```sh
dbreport run --config report.yml
```

Generate a report and override the output path:

```sh
dbreport run --config report.yml --output reports/daily.html
```

Send the generated report by email:

```sh
export DBREPORT_SMTP_USER='smtp-user'
export DBREPORT_SMTP_PASSWORD='smtp-password'
dbreport run --config report.yml --email
```

## Configuration

Primary config file name:

```text
report.yml
```

Example:

```yaml
title: "Daily MariaDB Summary"

output:
  html: "report.html"

database:
  host: "127.0.0.1"
  port: 3306
  name: "appdb"
  user_env: "DBREPORT_DB_USER"
  password_env: "DBREPORT_DB_PASSWORD"
  timeout_seconds: 10
  tls: false

limits:
  max_rows_per_query: 1000

email:
  enabled: false
  smtp_host: "smtp.example.com"
  smtp_port: 587
  starttls: true
  username_env: "DBREPORT_SMTP_USER"
  password_env: "DBREPORT_SMTP_PASSWORD"
  from: "reports@example.com"
  to:
    - "ops@example.com"
  subject: "Daily MariaDB Summary"
  send_html_body: true
  attach_html: false

queries:
  - id: "total_orders"
    title: "Total Orders"
    type: "metric"
    sql: |
      SELECT COUNT(*) AS value
      FROM orders
```

See `docs/CONFIGURATION.md` for the full schema.

## Report section types

Supported section types:

```text
metric
table
bar
line
pie
```

See `docs/REPORT_SECTIONS.md` for contracts and query shape examples.

## Email

Email support uses Go standard-library SMTP/TLS functionality. It can send the
HTML report as the email body, attach the generated HTML file, or both.

See `docs/EMAIL.md` for SMTP config details and operational notes.

## Security notes

Recommended defaults:

- Use a read-only MariaDB account.
- Use least-privilege grants for only the reporting schema/tables needed.
- Store DB and SMTP secrets in environment variables.
- Do not put secrets in `report.yml`.
- Avoid reporting sensitive fields unless the report recipients are approved to
  receive them.
- Use TLS/STARTTLS where available.
- Be careful when emailing reports outside the organization.

See `docs/SECURITY.md` for more detail.

## Optional MariaDB integration smoke test

Run the optional integration smoke test to validate end-to-end report generation against a temporary MariaDB instance:

```sh
./scripts/integration_mariadb_smoke.sh
```

The script prefers a local MariaDB runtime first, falls back to Docker, and otherwise exits with a clear skip message.

Run the smoke test with `DBREPORT_KEEP_SAMPLE_REPORT=1` to save a real generated sample report to `docs/assets/sample-report.html`. The sample data is fake and uses non-real domains with reserved TLDs such as `.invalid`, `.test`, and `.example`.

A generated sample HTML report is included at:

- [docs/assets/sample-report.html](docs/assets/sample-report.html)

It is generated from the optional MariaDB smoke-test dataset using fake data and reserved example domains.


The project is preparing for `v0.1.0-alpha.1` with validated smoke-test sample assets and release metadata.

## Build

```sh
go build -o dbreport ./cmd/dbreport
```

Release-style build with metadata:

```sh
go build \
  -ldflags="-X 'github.com/rswestmoreland/dbreport/internal/version.Version=0.1.0' -X 'github.com/rswestmoreland/dbreport/internal/version.Commit=$(git rev-parse --short HEAD)' -X 'github.com/rswestmoreland/dbreport/internal/version.Date=$(date -u +%Y-%m-%d)'" \
  -o dbreport ./cmd/dbreport
```

See `docs/BUILD_RELEASE.md` and `docs/RELEASE.md` for release packaging notes.

## Development validation

Run in a normal Go environment with network access to download modules:

```sh
go mod tidy
gofmt -w ./cmd ./internal
go test ./...
go build ./cmd/dbreport
```

The first `go mod tidy` run will create or update `go.sum`.

## License

MIT License.

Copyright (c) 2026 Richard S. Westmoreland

## Author

Richard S. Westmoreland  
dev@rswestmore.land

## Release packaging

Release scripts are included under `scripts/`.

Validate a release candidate:

```sh
./scripts/check_release.sh
```

Build release artifacts:

```sh
VERSION=0.1.0 \
COMMIT=$(git rev-parse --short HEAD) \
DATE=$(date -u +%Y-%m-%d) \
./scripts/build_release.sh
```


- sample-report.html is the authoritative generated visual sample.
- Report HTML head includes generator metadata and a project help link to https://github.com/rswestmoreland/dbreport; the report remains self-contained.
