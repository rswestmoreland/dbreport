# Roadmap

`dbreport` is intentionally small. The near-term goal is a practical
single-binary MariaDB reporting utility, not a full BI platform.

## Current scope

- YAML report configuration.
- MariaDB query execution.
- Metric, table, bar, line, and pie sections.
- Self-contained HTML output.
- Optional SMTP email delivery.
- Release packaging scripts.

## Candidate near-term hardening

- Keep smoke-tested sample report assets up to date.
- Add more example reports.
- Improve chart sizing for very large result sets.
- Add stricter config diagnostics where useful.
- Consider optional CI execution of MariaDB smoke tests as non-blocking.

## Future work: query parameter contract (not implemented yet)

Planned SQL placeholder style:

```sql
WHERE login_time >= :start_time
  AND login_time < :end_time
  AND result = :result
```

Planned CLI invocation with repeated parameters:

```sh
dbreport run --config auth-report.yml \
  --param start_time=2026-02-01T00:00:00Z \
  --param end_time=2026-02-08T00:00:00Z \
  --param result=success
```

Planned optional parameter file:

```sh
dbreport run --config auth-report.yml --params params.yml
```

Planned `params.yml` format:

```yaml
start_time: "2026-02-01T00:00:00Z"
end_time: "2026-02-08T00:00:00Z"
result: "success"
```

Rules for the planned contract:

- No parameter declaration block in `report.yml`.
- No type/required/default schema initially.
- If a query references `:name`, that parameter must be provided by `--param`
  or `--params`.
- Missing referenced parameters must fail with a clear error.
- No parameter defaults.
- Parameters are value-only.
- Parameters must be bound through `database/sql` after converting named
  placeholders to driver `?` placeholders.
- No raw string interpolation.
- Parameters must not control table names, column names, SQL keywords,
  `ORDER BY` expressions, SQL fragments, or arbitrary SQL structure.
- If dynamic identifiers are needed later, use explicit allowlisted config.

## Deferred features

- PDF output.
- Multiple databases per config.
- Multiple reports per config.
- Scheduled mode.
- Embedded web preview mode.
- Additional chart types.
- Theming.
- Redaction rules.
- CSV or JSON export.
- Vault integration.
