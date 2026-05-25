# Integration smoke test (optional)

Use `scripts/integration_mariadb_smoke.sh` to run an end-to-end MariaDB smoke test outside normal unit tests.

## What it validates

- Builds `dbreport` from source.
- Creates and seeds an authentication-reporting sample database.
- Runs `dbreport check` and `dbreport run` using `examples/auth-login-report.yml`.
- Verifies expected report titles/sections/chart SVG markup in generated HTML.
- Rejects external/scripted asset references in report output.

## Runtime strategy

The script tries the following in order:

1. Docker (`mariadb:11`) temporary container.
2. Local MariaDB/MySQL client+reachable server.
3. Clear skip exit when neither is available.

The smoke test is optional and is **not** part of default `go test ./...` or `scripts/check_release.sh`.

## Sample data and privacy

- Schema: `examples/auth-login-schema.sql`
- Seed data: `examples/auth-login-seed.sql`
- Report config: `examples/auth-login-report.yml`

Seed data is deterministic and uses fictional names plus non-real domains only:

- `example.invalid`
- `test.invalid`
- `corp.invalid`

No real customer or personal production data is used.

## Cleanup behavior

- Docker container is removed automatically on exit.
- Temporary workspace `/tmp/dbreport-smoke` (binary/config/report) is removed automatically.
- No persistent Docker volumes are created by the script.

## Saving a committed sample report

By default, the smoke test cleans up temporary report output under `/tmp/dbreport-smoke`.

To regenerate and save the documentation sample report:

```sh
DBREPORT_KEEP_SAMPLE_REPORT=1 ./scripts/integration_mariadb_smoke.sh
```

When `DBREPORT_KEEP_SAMPLE_REPORT=1` is set, the script copies the generated HTML report to:

- `docs/assets/sample-report.html`

The script still removes temporary MariaDB data directories, socket/pid/log files, temp binary files, and temp workspace content.

## Visual sample asset

The repository currently uses `docs/assets/sample-report.svg` as the README visual. A PNG visual can be generated from `docs/assets/sample-report.html` in environments with a headless browser toolchain.

If PNG regeneration is not available in your local environment, keep `docs/assets/sample-report.html` as the source of truth and update `docs/assets/sample-report.svg` as fallback documentation art.
