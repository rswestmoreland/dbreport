# dbreport Commands

## Overview

`dbreport` is intentionally small. Commands are focused on generating reports,
checking configuration, and showing metadata.

## Help

```sh
dbreport --help
dbreport help
```

Shows usage-focused help. Author and license details are intentionally kept in
`dbreport about` to keep help output concise.

## Version

```sh
dbreport version
```

Shows concise version and build metadata.

Example:

```text
dbreport 0.1.0
Build date: 2026-05-25
Commit: abc1234
```

## About

```sh
dbreport about
```

Shows project description, author, email, license, copyright, version, build date,
and commit.

## Sample config

```sh
dbreport sample-config > report.yml
```

Prints a usable YAML configuration sample.

## Check

```sh
dbreport check --config report.yml
```

Validates the YAML config, required environment variables, MariaDB connectivity,
and configured query execution.

Supported options:

```text
--config FILE
--email
--no-email
--quiet
--verbose
```

Notes:

- `--output` is not valid for `check` because check does not write a report.
- `--email` validates email configuration in addition to normal checks.
- `--quiet` and `--verbose` cannot be used together.
- `--email` and `--no-email` cannot be used together.

## Run

```sh
dbreport run --config report.yml
```

Runs configured queries and writes the configured HTML report.

Supported options:

```text
--config FILE
--output FILE
--email
--no-email
--quiet
--verbose
```

Examples:

```sh
dbreport run --config report.yml
dbreport run --config report.yml --output reports/daily.html
dbreport run --config report.yml --email
dbreport run --config report.yml --no-email
```

## Exit codes

```text
0  success
1  general error
2  config error
3  database error
4  query error
5  output/write error
6  email error
```
