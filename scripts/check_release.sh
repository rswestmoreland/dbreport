\
#!/usr/bin/env sh
set -eu

echo "running gofmt check"
unformatted="$(gofmt -l ./cmd ./internal 2>/dev/null || true)"
if [ -n "$unformatted" ]; then
    echo "gofmt required for:"
    echo "$unformatted"
    exit 1
fi

echo "running go mod tidy check"
cp go.mod /tmp/dbreport_go_mod_before.$$
cp go.sum /tmp/dbreport_go_sum_before.$$ 2>/dev/null || true
go mod tidy
if ! cmp -s go.mod /tmp/dbreport_go_mod_before.$$; then
    echo "go.mod changed after go mod tidy"
    rm -f /tmp/dbreport_go_mod_before.$$ /tmp/dbreport_go_sum_before.$$
    exit 1
fi
if [ -f /tmp/dbreport_go_sum_before.$$ ] && ! cmp -s go.sum /tmp/dbreport_go_sum_before.$$; then
    echo "go.sum changed after go mod tidy"
    rm -f /tmp/dbreport_go_mod_before.$$ /tmp/dbreport_go_sum_before.$$
    exit 1
fi
rm -f /tmp/dbreport_go_mod_before.$$ /tmp/dbreport_go_sum_before.$$

echo "running tests"
go test ./...

echo "running build"
go build ./cmd/dbreport

echo "release checks passed"
