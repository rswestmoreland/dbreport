# Release Packaging

This document describes the release process for `dbreport`.

## Release goals

A release should provide small, self-contained command-line binaries with version
metadata embedded at build time.

## Release artifact targets

Recommended first release targets:

```text
linux/amd64
linux/arm64
darwin/amd64
darwin/arm64
windows/amd64
```

Linux binaries are the primary target. Other targets are convenient because Go
cross-compilation is straightforward for this project.

## Required validation before release

Run these commands from the repository root:

```sh
gofmt -w ./cmd ./internal
go mod tidy
go test ./...
go build ./cmd/dbreport
```

Then run:

```sh
./scripts/check_release.sh
```

## Build release artifacts

Example:

```sh
VERSION=0.1.0 \
COMMIT=$(git rev-parse --short HEAD) \
DATE=$(date -u +%Y-%m-%d) \
./scripts/build_release.sh
```

Artifacts are written to:

```text
dist/
```

The release script also writes `dist/SHA256SUMS` when `sha256sum` or `shasum`
is available.

## Version metadata

The following values are embedded into the binary:

```text
internal/version.Version
internal/version.Commit
internal/version.Date
```

Verify the built binary:

```sh
./dist/dbreport_0.1.0_linux_amd64 version
./dist/dbreport_0.1.0_linux_amd64 about
```

## Release checklist

- [ ] `go mod tidy` completed.
- [ ] `go test ./...` passed.
- [ ] `go build ./cmd/dbreport` passed.
- [ ] `./scripts/check_release.sh` passed.
- [ ] `VERSION` set to the release version.
- [ ] `COMMIT` set to the release commit.
- [ ] `DATE` set to the UTC release date.
- [ ] `./scripts/build_release.sh` completed.
- [ ] `dist/SHA256SUMS` generated.
- [ ] `dbreport version` shows the release version.
- [ ] `dbreport about` shows author, license, version, date, and commit.
- [ ] `README.md` reflects current behavior.
- [ ] `CHANGELOG.md` has a release entry.
- [ ] `LICENSE` is included in the repository.
- [ ] Example configs are included.
