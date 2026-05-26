# dbreport Configuration

## File format

Primary config file:

```text
report.yml
```

YAML is parsed with `github.com/goccy/go-yaml`, compiled into the Go binary.
This does not add an external runtime dependency.

## YAML feature guardrail

Use simple YAML features only:

- mappings
- sequences
- strings
- numbers
- booleans
- block scalars for SQL

Avoid relying on:

- anchors
- aliases
- merge keys
- custom tags

This keeps report files readable and portable.

## Top-level schema

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

queries: []
```

## `title`

Required report title shown in the HTML report and used as a default email subject
when appropriate.

## `output.html`

Required path for the generated HTML report.

Example:

```yaml
output:
  html: "reports/daily.html"
```

Parent directories are created when needed.

## `database`

Required fields:

```yaml
database:
  host: "127.0.0.1"
  port: 3306
  name: "appdb"
  user_env: "DBREPORT_DB_USER"
  password_env: "DBREPORT_DB_PASSWORD"
  timeout_seconds: 10
  tls: false
```

`user_env` and `password_env` name environment variables. The actual secret values
must not be stored in the config file.

## `limits.max_rows_per_query`

Required. Prevents accidental huge report sections.

Example:

```yaml
limits:
  max_rows_per_query: 1000
```

## `email`

Optional unless email is enabled in the config or `--email` is passed.

See `docs/EMAIL.md`.

## `queries`

At least one query is required.

Common fields:

```yaml
queries:
  - id: "recent_orders"
    title: "Recent Orders"
    type: "table"
    sql: |
      SELECT id, status, created_at
      FROM orders
      ORDER BY created_at DESC
      LIMIT 25
```

Rules:

- `id` must be unique.
- `title` is required.
- `type` must be one of `metric`, `table`, `bar`, `line`, or `pie`.
- `sql` must not be empty.
- Chart sections require `label_column` and `value_column`. For `line`, optional `series_column` enables multi-series lines. Optional `show_table` (default true) hides fallback tables when false for bar/line/pie.

## Environment variables

Example:

```sh
export DBREPORT_DB_USER='report_user'
export DBREPORT_DB_PASSWORD='change-me'
export DBREPORT_SMTP_USER='smtp-user'
export DBREPORT_SMTP_PASSWORD='smtp-password'
```


- sample-report.html is the authoritative generated visual sample.


## Safety defaults and overrides

- `allowed_databases` and `allowed_tables` are defense-in-depth controls.
- If `blocked_functions`, `blocked_columns`, or `blocked_patterns` are omitted, built-in defaults are used.
- If provided, user lists override defaults in this release (no merge).
- Static SQL validation is heuristic and does not replace a read-only DB account.
- A single `report.yml` is acceptable when it contains no secrets.
- Future optional split report/db-auth config may be added later.

## Query parameters
Use repeated --param name=value and optional --params params.yml. CLI --param overrides file values. Duplicate --param names are rejected. params.yml must be a top-level mapping of scalar values only.
