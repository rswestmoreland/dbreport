#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORK_DIR="/tmp/dbreport-smoke"
DB_NAME="dbreport_smoke"
DB_USER="dbreport_smoke_user"
DB_PASSWORD="smoke_pw_local_only"
DB_HOST="127.0.0.1"
DB_PORT="3306"
CONTAINER_NAME="dbreport-smoke-mariadb-$RANDOM"
MARIADB_IMAGE="mariadb:11"
REPORT_PATH="$WORK_DIR/report.html"

SCHEMA_FILE="$ROOT_DIR/examples/auth-login-schema.sql"
SEED_FILE="$ROOT_DIR/examples/auth-login-seed.sql"
REPORT_CONFIG="$WORK_DIR/auth-login-report.yml"

MODE=""
SKIP_REASONS=()
CONTAINER_STARTED=0

cleanup() {
  if [[ "$CONTAINER_STARTED" -eq 1 ]]; then
    docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true
  fi
  rm -rf "$WORK_DIR"
}
trap cleanup EXIT

mkdir -p "$WORK_DIR"
cp "$ROOT_DIR/examples/auth-login-report.yml" "$REPORT_CONFIG"

log() { printf '%s\n' "$*"; }
skip() { log "SKIP: $*"; exit 2; }

if command -v docker >/dev/null 2>&1; then
  MODE="docker"
elif command -v mariadb >/dev/null 2>&1 || command -v mysql >/dev/null 2>&1; then
  MODE="local"
else
  SKIP_REASONS+=("docker not found")
  SKIP_REASONS+=("mariadb/mysql client not found")
  skip "No MariaDB runtime available (${SKIP_REASONS[*]})."
fi

MYSQL_CLIENT="mariadb"
if ! command -v mariadb >/dev/null 2>&1; then
  MYSQL_CLIENT="mysql"
fi

if [[ "$MODE" == "docker" ]]; then
  log "Using Docker-based MariaDB smoke test"
  docker run -d --rm \
    --name "$CONTAINER_NAME" \
    -e MARIADB_ROOT_PASSWORD="$DB_PASSWORD" \
    -e MARIADB_DATABASE="$DB_NAME" \
    -e MARIADB_USER="$DB_USER" \
    -e MARIADB_PASSWORD="$DB_PASSWORD" \
    -p "$DB_PORT:3306" \
    "$MARIADB_IMAGE" >/dev/null
  CONTAINER_STARTED=1

  ready=0
  for _ in $(seq 1 60); do
    if docker exec "$CONTAINER_NAME" mariadb -uroot -p"$DB_PASSWORD" -e 'SELECT 1' >/dev/null 2>&1; then
      ready=1
      break
    fi
    sleep 1
  done
  [[ "$ready" -eq 1 ]] || { log "MariaDB container failed readiness check"; exit 1; }

  docker exec -i "$CONTAINER_NAME" mariadb -uroot -p"$DB_PASSWORD" "$DB_NAME" < "$SCHEMA_FILE"
  docker exec -i "$CONTAINER_NAME" mariadb -uroot -p"$DB_PASSWORD" "$DB_NAME" < "$SEED_FILE"
else
  log "Using local MariaDB client/server smoke test"
  if ! command -v "$MYSQL_CLIENT" >/dev/null 2>&1; then
    skip "Local mode selected but no mariadb/mysql client available"
  fi

  export DBREPORT_DB_USER="${DBREPORT_DB_USER:-$DB_USER}"
  export DBREPORT_DB_PASSWORD="${DBREPORT_DB_PASSWORD:-$DB_PASSWORD}"
  DB_USER="$DBREPORT_DB_USER"
  DB_PASSWORD="$DBREPORT_DB_PASSWORD"

  if ! "$MYSQL_CLIENT" -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e 'SELECT 1' >/dev/null 2>&1; then
    skip "Local MariaDB is not reachable at $DB_HOST:$DB_PORT for user $DB_USER"
  fi
  "$MYSQL_CLIENT" -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME" >/dev/null
  "$MYSQL_CLIENT" -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$SCHEMA_FILE"
  "$MYSQL_CLIENT" -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$SEED_FILE"
fi

if [[ "$MODE" == "docker" ]]; then
  export DBREPORT_DB_USER="$DB_USER"
  export DBREPORT_DB_PASSWORD="$DB_PASSWORD"
fi

cd "$ROOT_DIR"
go build -o "$WORK_DIR/dbreport" ./cmd/dbreport

"$WORK_DIR/dbreport" check --config "$REPORT_CONFIG"
"$WORK_DIR/dbreport" run --config "$REPORT_CONFIG"

[[ -f "$REPORT_PATH" ]] || { log "report file missing: $REPORT_PATH"; exit 1; }
[[ -s "$REPORT_PATH" ]] || { log "report file empty: $REPORT_PATH"; exit 1; }

grep -q "Authentication Login Activity Report" "$REPORT_PATH"
grep -q "Last Logins by User" "$REPORT_PATH"
grep -q "Successful vs Failed Logins" "$REPORT_PATH"
grep -q "Top Browsers Used" "$REPORT_PATH"
grep -q "Top Countries Logged in From" "$REPORT_PATH"
grep -q "<svg" "$REPORT_PATH"

if grep -Eqi 'http://|https://|<script|\bsrc=|\bhref=' "$REPORT_PATH"; then
  log "report contains disallowed external/script references"
  exit 1
fi

log "MariaDB integration smoke test passed ($MODE mode)."
