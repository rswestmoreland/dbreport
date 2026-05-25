#!/usr/bin/env sh
set -eu

APP_NAME="dbreport"
MODULE_PATH="github.com/rswestmoreland/dbreport"
VERSION="${VERSION:-0.1.0-dev}"
COMMIT="${COMMIT:-unknown}"
DATE="${DATE:-$(date -u +%Y-%m-%d)}"
DIST_DIR="${DIST_DIR:-dist}"

mkdir -p "$DIST_DIR"

build_one() {
    goos="$1"
    goarch="$2"
    suffix="$3"

    output="$DIST_DIR/${APP_NAME}_${VERSION}_${goos}_${goarch}${suffix}"

    echo "building $output"
    GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build \
        -trimpath \
        -ldflags="-s -w -X '${MODULE_PATH}/internal/version.Version=${VERSION}' -X '${MODULE_PATH}/internal/version.Commit=${COMMIT}' -X '${MODULE_PATH}/internal/version.Date=${DATE}'" \
        -o "$output" ./cmd/dbreport
}

build_one linux amd64 ""
build_one linux arm64 ""
build_one darwin amd64 ""
build_one darwin arm64 ""
build_one windows amd64 ".exe"

if command -v sha256sum >/dev/null 2>&1; then
    (
        cd "$DIST_DIR"
        sha256sum * > SHA256SUMS
    )
elif command -v shasum >/dev/null 2>&1; then
    (
        cd "$DIST_DIR"
        shasum -a 256 * > SHA256SUMS
    )
else
    echo "warning: neither sha256sum nor shasum was found; checksums were not generated" >&2
fi

echo "release artifacts written to $DIST_DIR"
