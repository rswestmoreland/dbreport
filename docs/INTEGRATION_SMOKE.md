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
