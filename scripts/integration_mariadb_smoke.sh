#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORK_DIR="/tmp/dbreport-smoke"
DB_NAME="dbreport_test"
DB_USER="dbreport_smoke_user"
DB_PASSWORD="smoke_pw_local_only"
DB_HOST="127.0.0.1"
DB_PORT="13306"
DB_SOCKET="$WORK_DIR/mariadb.sock"
DB_PID_FILE="$WORK_DIR/mariadb.pid"
DB_LOG_FILE="$WORK_DIR/mariadb.log"
DB_DATA_DIR="$WORK_DIR/mariadb-data"
CONTAINER_NAME="dbreport-smoke-mariadb-$RANDOM"
MARIADB_IMAGE="mariadb:11"
REPORT_PATH="$WORK_DIR/report.html"
REPORT_OUTPUT_DIR="/tmp/dbreport-smoke"
KEEP_SAMPLE_REPORT="${DBREPORT_KEEP_SAMPLE_REPORT:-0}"
DOC_SAMPLE_REPORT="$ROOT_DIR/docs/assets/sample-report.html"

SCHEMA_FILE="$ROOT_DIR/examples/auth-login-schema.sql"
SEED_FILE="$ROOT_DIR/examples/auth-login-seed.sql"
REPORT_CONFIG="$WORK_DIR/auth-login-report.yml"

MODE=""
CONTAINER_STARTED=0
LOCAL_SERVER_STARTED=0

cleanup() {
  if [[ "$LOCAL_SERVER_STARTED" -eq 1 && -f "$DB_PID_FILE" ]]; then
    local pid
    pid="$(cat "$DB_PID_FILE" 2>/dev/null || true)"
    if [[ -n "$pid" ]] && kill -0 "$pid" >/dev/null 2>&1; then
      kill "$pid" >/dev/null 2>&1 || true
      for _ in $(seq 1 30); do
        if ! kill -0 "$pid" >/dev/null 2>&1; then
          break
        fi
        sleep 0.2
      done
      kill -9 "$pid" >/dev/null 2>&1 || true
    fi
  fi
  if [[ "$CONTAINER_STARTED" -eq 1 ]]; then
    docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true
  fi
  rm -rf "$WORK_DIR"
  rm -f "$REPORT_OUTPUT_DIR/report.html"
}
trap cleanup EXIT

mkdir -p "$WORK_DIR"
cp "$ROOT_DIR/examples/auth-login-report.yml" "$REPORT_CONFIG"
sed -i "s/^  port: .*/  port: $DB_PORT/" "$REPORT_CONFIG"
sed -i "s/^  name: .*/  name: \"$DB_NAME\"/" "$REPORT_CONFIG"

log() { printf '%s\n' "$*"; }
skip() { log "SKIP: $*"; exit 2; }

MYSQL_CLIENT="mariadb"
if ! command -v mariadb >/dev/null 2>&1; then
  MYSQL_CLIENT="mysql"
fi

MARIADB_SERVER_BIN=""
if command -v mariadbd >/dev/null 2>&1; then
  MARIADB_SERVER_BIN="$(command -v mariadbd)"
elif command -v mysqld >/dev/null 2>&1; then
  MARIADB_SERVER_BIN="$(command -v mysqld)"
fi

if [[ -n "$MARIADB_SERVER_BIN" ]] && command -v "$MYSQL_CLIENT" >/dev/null 2>&1; then
  MODE="local"
elif command -v docker >/dev/null 2>&1; then
  MODE="docker"
else
  skip "No MariaDB runtime available (need local mariadbd/mysqld + mariadb/mysql client, or Docker)."
fi

start_local_mariadb() {
  log "Using local temporary MariaDB smoke test"

  rm -rf "$DB_DATA_DIR"
  mkdir -p "$DB_DATA_DIR"

  if command -v mariadb-install-db >/dev/null 2>&1; then
    mariadb-install-db --datadir="$DB_DATA_DIR" --auth-root-authentication-method=normal --skip-test-db >/dev/null
  elif command -v mysql_install_db >/dev/null 2>&1; then
    mysql_install_db --datadir="$DB_DATA_DIR" --auth-root-authentication-method=normal --skip-test-db >/dev/null
  else
    skip "Unable to initialize local MariaDB datadir (no mariadb-install-db/mysql_install_db)"
  fi

  "$MARIADB_SERVER_BIN" \
    --datadir="$DB_DATA_DIR" \
    --socket="$DB_SOCKET" \
    --pid-file="$DB_PID_FILE" \
    --port="$DB_PORT" \
    --bind-address=127.0.0.1 \
    --skip-networking=0 \
    --skip-grant-tables \
    --log-error="$DB_LOG_FILE" \
    --user=root >/dev/null 2>&1 &

  LOCAL_SERVER_STARTED=1

  local ready=0
  for _ in $(seq 1 90); do
    if "$MYSQL_CLIENT" --protocol=TCP -h "$DB_HOST" -P "$DB_PORT" -u root -e 'SELECT 1' >/dev/null 2>&1; then
      ready=1
      break
    fi
    sleep 1
  done

  if [[ "$ready" -ne 1 ]]; then
    log "Local MariaDB failed readiness check"
    [[ -f "$DB_LOG_FILE" ]] && tail -n 120 "$DB_LOG_FILE" || true
    exit 1
  fi

  "$MYSQL_CLIENT" --protocol=TCP -h "$DB_HOST" -P "$DB_PORT" -u root -e "CREATE DATABASE IF NOT EXISTS $DB_NAME" >/dev/null
  "$MYSQL_CLIENT" --protocol=TCP -h "$DB_HOST" -P "$DB_PORT" -u root "$DB_NAME" < "$SCHEMA_FILE"
  "$MYSQL_CLIENT" --protocol=TCP -h "$DB_HOST" -P "$DB_PORT" -u root "$DB_NAME" < "$SEED_FILE"

  export DBREPORT_DB_USER="root"
  export DBREPORT_DB_PASSWORD="$DB_PASSWORD"
}

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

  export DBREPORT_DB_USER="$DB_USER"
  export DBREPORT_DB_PASSWORD="$DB_PASSWORD"
else
  start_local_mariadb
fi

cd "$ROOT_DIR"
SAMPLE_VERSION="0.1.0-alpha.1"
SAMPLE_COMMIT="$(git rev-parse --short HEAD)"
SAMPLE_DATE="$(date -u +%Y-%m-%d)"
go build -ldflags="-X 'github.com/rswestmoreland/dbreport/internal/version.Version=${SAMPLE_VERSION}' -X 'github.com/rswestmoreland/dbreport/internal/version.Commit=${SAMPLE_COMMIT}' -X 'github.com/rswestmoreland/dbreport/internal/version.Date=${SAMPLE_DATE}'" -o "$WORK_DIR/dbreport" ./cmd/dbreport

"$WORK_DIR/dbreport" check --config "$REPORT_CONFIG"
"$WORK_DIR/dbreport" run --config "$REPORT_CONFIG" --output "$REPORT_PATH"

[[ -f "$REPORT_PATH" ]] || { log "report file missing: $REPORT_PATH"; exit 1; }
[[ -s "$REPORT_PATH" ]] || { log "report file empty: $REPORT_PATH"; exit 1; }

grep -q "Authentication Login Activity Report" "$REPORT_PATH"
grep -q "Last Logins by User" "$REPORT_PATH"
grep -q "Successful vs Failed Logins" "$REPORT_PATH"
grep -q "Top Browsers Used" "$REPORT_PATH"
grep -q "Top Countries Logged in From" "$REPORT_PATH"
grep -q "<svg" "$REPORT_PATH"

if grep -Eqi '<script|\bsrc=|rel="stylesheet"|url\(' "$REPORT_PATH"; then
  log "report contains disallowed external/script references"
  exit 1
fi

if grep -Fq 'https://github.com/rswestmoreland/dbreport' "$REPORT_PATH"; then
  true
else
  log "report missing expected project metadata URL"
  exit 1
fi

if [[ "$KEEP_SAMPLE_REPORT" == "1" ]]; then
  mkdir -p "$(dirname "$DOC_SAMPLE_REPORT")"
  cp "$REPORT_PATH" "$DOC_SAMPLE_REPORT"
  log "Saved sample report to $DOC_SAMPLE_REPORT"
fi

log "MariaDB integration smoke test passed ($MODE mode)."
