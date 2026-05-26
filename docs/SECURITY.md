# dbreport Security Notes

`dbreport` is intentionally simple, but it still handles database access and may
produce reports containing sensitive data.

## Database account

Use a read-only MariaDB user.

Example grant shape:

```sql
CREATE USER 'report_user'@'%' IDENTIFIED BY 'change-me';
GRANT SELECT ON appdb.* TO 'report_user'@'%';
FLUSH PRIVILEGES;
```

Restrict grants further when possible.

## Secrets

Do not store secrets in `report.yml`.

Use environment variables:

```sh
export DBREPORT_DB_USER='report_user'
export DBREPORT_DB_PASSWORD='change-me'
```

The config file should contain only environment variable names:

```yaml
database:
  user_env: "DBREPORT_DB_USER"
  password_env: "DBREPORT_DB_PASSWORD"
```

## Query safety

MVP guardrails:

- SQL comes from the config file, not command-line arguments.
- Query timeout is required.
- Row cap is required.
- The tool expects a read-only DB account.

Recommended query practices:

- Use explicit column lists.
- Use `LIMIT` for detail tables.
- Avoid `SELECT *` in reports.
- Avoid sensitive fields unless required.
- Keep expensive reports indexed and tested.

## Email safety

Reports may contain sensitive data. Before enabling email:

- Review report contents.
- Confirm recipients.
- Use STARTTLS where available.
- Avoid external recipients unless approved.
- Use a dedicated SMTP account where practical.

## Logging

The tool should not print passwords, full credential-bearing DSNs, or secret
environment variable values.


Safety defaults: omitted blocked function/column/pattern lists use built-in defaults; configured lists override defaults for that list.

Named parameters are bound through database/sql arguments, not SQL interpolation. Structural placeholders (table/column/order/group/sql fragments) are rejected.
