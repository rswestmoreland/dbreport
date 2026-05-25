# Roadmap

`dbreport` is intentionally small. The near-term goal is a practical
single-binary MariaDB reporting utility, not a full BI platform.

## Current scope

- YAML report configuration.
- MariaDB query execution.
- Metric, table, bar chart, and line chart sections.
- Self-contained HTML output.
- Optional SMTP email delivery.
- Release packaging scripts.

## Candidate near-term hardening

- Complete external build and test validation.
- Add integration tests using a temporary MariaDB test instance.
- Add more example reports.
- Improve chart sizing for very large result sets.
- Add stricter config diagnostics where useful.
- Add optional `EXPLAIN` validation mode for queries.

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
