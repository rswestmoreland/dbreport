# dbreport Build and Release

## Development build

```sh
go build -o dbreport ./cmd/dbreport
```

## Validation

Run in a normal Go environment with network access to download modules:

```sh
go mod tidy
go fmt ./...
go test ./...
go build ./cmd/dbreport
```

## Release build

```sh
go build \
  -ldflags="-X 'github.com/rswestmoreland/dbreport/internal/version.Version=0.1.0' -X 'github.com/rswestmoreland/dbreport/internal/version.Commit=$(git rev-parse --short HEAD)' -X 'github.com/rswestmoreland/dbreport/internal/version.Date=$(date -u +%Y-%m-%d)'" \
  -o dbreport ./cmd/dbreport
```

## Suggested release artifacts

- `dbreport` Linux amd64 binary.
- Optional Linux arm64 binary.
- `LICENSE`.
- `README.md`.
- `examples/report.yml`.
- Checksums.
- Changelog entry.

## Versioning

Start with:

```text
0.1.0
```

Use semantic versioning for public releases.
