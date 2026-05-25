APP_NAME := dbreport
VERSION ?= 0.1.0-dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%d)

.PHONY: help fmt test build release clean check

help:
	@echo "Targets:"
	@echo "  fmt      Format Go source"
	@echo "  test     Run Go tests"
	@echo "  build    Build local binary"
	@echo "  release  Build release artifacts into dist/"
	@echo "  check    Run release checks"
	@echo "  clean    Remove build artifacts"

fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

build:
	go build -trimpath \
		-ldflags="-X 'github.com/rswestmoreland/dbreport/internal/version.Version=$(VERSION)' -X 'github.com/rswestmoreland/dbreport/internal/version.Commit=$(COMMIT)' -X 'github.com/rswestmoreland/dbreport/internal/version.Date=$(DATE)'" \
		-o $(APP_NAME) ./cmd/dbreport

release:
	VERSION="$(VERSION)" COMMIT="$(COMMIT)" DATE="$(DATE)" ./scripts/build_release.sh

check:
	./scripts/check_release.sh

clean:
	rm -f $(APP_NAME)
	rm -rf dist build release
