# Security Model

dbreport assumes trusted, reviewed YAML query configuration and enforces defense-in-depth controls.

- Threat model: prevent accidental unsafe SQL and unsafe HTML/email output in generated reports.
- No command-line arbitrary SQL execution is supported.
- Query policy permits only SELECT and WITH ... SELECT patterns.
- Query policy blocks dangerous statements and tokens and rejects multi-statement attempts.
- Safety controls support blocked functions, blocked column names/patterns, and optional allowed tables/databases.
- Use a read-only database account with SELECT-only grants.
- Keep secrets in environment variables, not in config.
- Errors are wrapped with query id/title and avoid credential-bearing DSN output.
- HTML output escapes dynamic content and includes CSP.
- Email uses the same generated safe HTML; emailed reports may still contain sensitive operational data.
- limits.max_rows_per_query and limits.max_cell_length reduce report abuse risk.
- Static SQL validation is heuristic and not a full parser.
- Version-control report config and review queries before enabling email delivery.
- A single report.yml is acceptable when it contains no secrets.
- A future optional split between report and auth config may be added.
