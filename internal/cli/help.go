package cli

import (
	"fmt"
	"io"

	"github.com/rswestmoreland/dbreport/internal/version"
)

func PrintHelp(w io.Writer) {
	fmt.Fprintf(w, `%s - %s

Usage:
  dbreport run --config report.yml [--email] [--quiet] [--verbose]
  dbreport check --config report.yml [--quiet] [--verbose]
  dbreport sample-config
  dbreport version
  dbreport about
  dbreport help

Commands:
  run            Run configured queries and generate a report
  check          Validate configuration and database connectivity
  sample-config  Print a sample configuration file
  version        Show concise version and build information
  about          Show project, author, license, and build information
  help           Show this help text

Options:
  --config FILE   Path to report configuration
  --email         Send generated report by email
  --no-email      Disable configured email for this run/check
  --output FILE   Override configured output path for run
  --quiet, -q     Suppress successful status output
  --verbose       Print additional progress details to stderr
  --help          Show help

Use "dbreport about" for author, license, and build information.
`, version.AppName, version.Description)
}

func printRunHelp(w io.Writer) {
	fmt.Fprint(w, `Usage:
  dbreport run --config report.yml [--output FILE] [--email] [--no-email] [--quiet] [--verbose]

Options:
  --config FILE   Path to report configuration. Defaults to report.yml.
  --output FILE   Override configured output HTML path.
  --email         Send generated report by email.
  --no-email      Disable email even when email.enabled is true in config.
  --quiet, -q     Suppress successful status output.
  --verbose       Print additional progress details to stderr.
  --help, -h      Show this help text.
`)
}

func printCheckHelp(w io.Writer) {
	fmt.Fprint(w, `Usage:
  dbreport check --config report.yml [--email] [--no-email] [--quiet] [--verbose]

Options:
  --config FILE   Path to report configuration. Defaults to report.yml.
  --email         Validate email settings even when email.enabled is false.
  --no-email      Skip email validation even when email.enabled is true in config.
  --quiet, -q     Suppress successful status output.
  --verbose       Print additional progress details to stderr.
  --help, -h      Show this help text.
`)
}
