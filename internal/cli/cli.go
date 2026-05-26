package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rswestmoreland/dbreport/internal/config"
	dbreportdb "github.com/rswestmoreland/dbreport/internal/db"
	dbreportemail "github.com/rswestmoreland/dbreport/internal/email"
	"github.com/rswestmoreland/dbreport/internal/output"
	"github.com/rswestmoreland/dbreport/internal/querypolicy"
	"github.com/rswestmoreland/dbreport/internal/report"
)

var defaultBlockedColumns = []string{"password", "passwd", "password_hash", "token", "api_key", "secret", "private_key", "session", "cookie"}

const (
	ExitSuccess       = 0
	ExitGeneralError  = 1
	ExitConfigError   = 2
	ExitDatabaseError = 3
	ExitQueryError    = 4
	ExitOutputError   = 5
	ExitEmailError    = 6
)

type Runner struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (r Runner) Run(args []string) int {
	if len(args) == 0 {
		PrintHelp(r.Stdout)
		return ExitSuccess
	}

	cmd := args[0]
	switch cmd {
	case "--help", "-h", "help":
		PrintHelp(r.Stdout)
		return ExitSuccess
	case "version", "--version", "-v":
		PrintVersion(r.Stdout)
		return ExitSuccess
	case "about":
		PrintAbout(r.Stdout)
		return ExitSuccess
	case "sample-config":
		return r.sampleConfig(args[1:])
	case "check":
		return r.check(args[1:])
	case "run":
		return r.run(args[1:])
	default:
		if strings.HasPrefix(cmd, "-") {
			fmt.Fprintf(r.Stderr, "unknown option: %s\n\n", cmd)
		} else {
			fmt.Fprintf(r.Stderr, "unknown command: %s\n\n", cmd)
		}
		PrintHelp(r.Stderr)
		return ExitGeneralError
	}
}

func (r Runner) sampleConfig(args []string) int {
	if len(args) > 0 {
		if isHelpArg(args[0]) {
			fmt.Fprintln(r.Stdout, "Usage: dbreport sample-config")
			return ExitSuccess
		}
		fmt.Fprintf(r.Stderr, "sample-config does not accept arguments: %s\n", strings.Join(args, " "))
		return ExitGeneralError
	}
	fmt.Fprint(r.Stdout, config.SampleYAML)
	return ExitSuccess
}

func (r Runner) check(args []string) int {
	opts, err := parseCommonOptions(args)
	if err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", err.Error())
		return ExitConfigError
	}
	if opts.Help {
		printCheckHelp(r.Stdout)
		return ExitSuccess
	}
	if opts.OutputPath != "" {
		fmt.Fprintln(r.Stderr, "--output is only valid with run")
		return ExitConfigError
	}

	cfg, path, code := r.loadConfig(opts)
	if code != ExitSuccess {
		return code
	}
	r.verbosef(opts, "config loaded: %s\n", path)
	r.verbosef(opts, "database target: %s\n", dbreportdb.SafeTarget(cfg.Database))

	handle, timeout, code := r.openDatabase(context.Background(), cfg)
	if code != ExitSuccess {
		return code
	}
	defer handle.Close()

	r.verbosef(opts, "running configured queries: %d\n", len(cfg.Queries))
	results, err := dbreportdb.RunAll(context.Background(), handle, cfg.Queries, cfg.Limits.MaxRowsPerQuery, timeout)
	if err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", err.Error())
		return ExitQueryError
	}

	blockedColumns := cfg.Safety.BlockedColumns
	if len(blockedColumns) == 0 {
		blockedColumns = defaultBlockedColumns
	}
	_, blockedPatterns := querypolicy.Defaults(cfg.Safety.BlockedFunctions, cfg.Safety.BlockedPatterns)
	for _, result := range results {
		if err := dbreportdb.ValidateReturnedColumns(result, blockedColumns, blockedPatterns); err != nil {
			fmt.Fprintf(r.Stderr, "%s\n", err.Error())
			return ExitQueryError
		}
	}

	r.statusf(opts, "configuration valid: %s\n", path)
	r.statusf(opts, "database connected: %s\n", dbreportdb.SafeTarget(cfg.Database))
	r.statusf(opts, "queries valid: %d\n", len(results))
	if !opts.Quiet {
		printQuerySummary(r.Stdout, results)
	}
	if opts.Email || cfg.Email.Enabled {
		r.statusf(opts, "email configuration valid\n")
	}
	return ExitSuccess
}

func (r Runner) run(args []string) int {
	opts, err := parseCommonOptions(args)
	if err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", err.Error())
		return ExitConfigError
	}
	if opts.Help {
		printRunHelp(r.Stdout)
		return ExitSuccess
	}
	cfg, path, code := r.loadConfig(opts)
	if code != ExitSuccess {
		return code
	}
	r.verbosef(opts, "config loaded: %s\n", path)
	if opts.OutputPath != "" {
		cfg.Output.HTML = opts.OutputPath
	}
	r.verbosef(opts, "database target: %s\n", dbreportdb.SafeTarget(cfg.Database))
	r.verbosef(opts, "report output: %s\n", cfg.Output.HTML)

	handle, timeout, code := r.openDatabase(context.Background(), cfg)
	if code != ExitSuccess {
		return code
	}
	defer handle.Close()

	r.verbosef(opts, "running configured queries: %d\n", len(cfg.Queries))
	results, err := dbreportdb.RunAll(context.Background(), handle, cfg.Queries, cfg.Limits.MaxRowsPerQuery, timeout)
	if err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", err.Error())
		return ExitQueryError
	}

	blockedColumns := cfg.Safety.BlockedColumns
	if len(blockedColumns) == 0 {
		blockedColumns = defaultBlockedColumns
	}
	_, blockedPatterns := querypolicy.Defaults(cfg.Safety.BlockedFunctions, cfg.Safety.BlockedPatterns)
	for _, result := range results {
		if err := dbreportdb.ValidateReturnedColumns(result, blockedColumns, blockedPatterns); err != nil {
			fmt.Fprintf(r.Stderr, "%s\n", err.Error())
			return ExitQueryError
		}
	}

	doc, err := report.NewDocument(*cfg, results)
	if err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", err.Error())
		return ExitOutputError
	}

	rendered, err := report.RenderHTML(doc)
	if err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", err.Error())
		return ExitOutputError
	}

	r.verbosef(opts, "writing report: %s\n", cfg.Output.HTML)
	if err := output.WriteFile(cfg.Output.HTML, rendered); err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", err.Error())
		return ExitOutputError
	}

	emailRequested := opts.Email || cfg.Email.Enabled
	if emailRequested {
		r.verbosef(opts, "sending report email to %d %s\n", len(cfg.Email.To), pluralize("recipient", len(cfg.Email.To)))
		if err := dbreportemail.Send(dbreportemail.SendRequest{
			Config:      cfg.Email,
			HTML:        rendered,
			ReportPath:  cfg.Output.HTML,
			GeneratedAt: doc.GeneratedAt,
		}); err != nil {
			fmt.Fprintf(r.Stderr, "%s\n", err.Error())
			fmt.Fprintf(r.Stderr, "report written: %s\n", cfg.Output.HTML)
			return ExitEmailError
		}
	}

	r.statusf(opts, "queries completed: %d\n", len(results))
	if !opts.Quiet {
		printQuerySummary(r.Stdout, results)
	}
	r.statusf(opts, "report written: %s\n", cfg.Output.HTML)
	if emailRequested {
		r.statusf(opts, "report emailed: %d %s\n", len(cfg.Email.To), pluralize("recipient", len(cfg.Email.To)))
	}
	return ExitSuccess
}

func (r Runner) loadConfig(opts commonOptions) (*config.Config, string, int) {
	path := opts.ConfigPath
	if path == "" {
		path = config.DefaultConfigPath
	}

	cfg, err := config.LoadFile(path)
	if err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", err.Error())
		return nil, path, ExitConfigError
	}
	if opts.NoEmail {
		cfg.Email.Enabled = false
		err = cfg.Validate()
	} else if opts.Email {
		err = cfg.ValidateForEmailRequest()
	} else {
		err = cfg.Validate()
	}
	if err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", err.Error())
		return nil, path, ExitConfigError
	}
	for _, q := range cfg.Queries {
		blockedFns, blockedPatterns := querypolicy.Defaults(cfg.Safety.BlockedFunctions, cfg.Safety.BlockedPatterns)
		err = querypolicy.Validate(q.ID, q.Title, q.SQL, querypolicy.SafetyOptions{ActiveDatabase: cfg.Database.Name, AllowedDatabases: cfg.Safety.AllowedDatabases, AllowedTables: cfg.Safety.AllowedTables, BlockedFunctions: blockedFns, BlockedPatterns: blockedPatterns})
		if err != nil {
			fmt.Fprintf(r.Stderr, "%s\n", err.Error())
			return nil, path, ExitConfigError
		}
	}
	return cfg, path, ExitSuccess
}

func (r Runner) openDatabase(ctx context.Context, cfg *config.Config) (*sql.DB, time.Duration, int) {
	timeout := time.Duration(cfg.Database.TimeoutSeconds) * time.Second
	handle, err := dbreportdb.Open(ctx, cfg.Database)
	if err != nil {
		fmt.Fprintf(r.Stderr, "%s\n", err.Error())
		return nil, timeout, ExitDatabaseError
	}
	return handle, timeout, ExitSuccess
}

func (r Runner) statusf(opts commonOptions, format string, args ...any) {
	if opts.Quiet {
		return
	}
	fmt.Fprintf(r.Stdout, format, args...)
}

func (r Runner) verbosef(opts commonOptions, format string, args ...any) {
	if !opts.Verbose || opts.Quiet {
		return
	}
	fmt.Fprintf(r.Stderr, format, args...)
}

func printQuerySummary(w io.Writer, results []dbreportdb.Result) {
	for _, result := range results {
		truncated := ""
		if result.Truncated {
			truncated = " (truncated at configured row cap)"
		}
		fmt.Fprintf(w, "- %s (%s): %d %s in %s%s\n",
			result.Query.ID,
			result.Query.Type,
			result.RowCount(),
			pluralize("row", result.RowCount()),
			formatDuration(result.Duration),
			truncated,
		)
	}
}

func formatDuration(value time.Duration) string {
	if value < time.Millisecond {
		return value.String()
	}
	return value.Round(time.Millisecond).String()
}

func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

type commonOptions struct {
	ConfigPath string
	OutputPath string
	Email      bool
	NoEmail    bool
	Quiet      bool
	Verbose    bool
	Help       bool
}

func parseCommonOptions(args []string) (commonOptions, error) {
	var opts commonOptions
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--help", "-h":
			opts.Help = true
		case "--email":
			opts.Email = true
		case "--no-email":
			opts.NoEmail = true
		case "--quiet", "-q":
			opts.Quiet = true
		case "--verbose":
			opts.Verbose = true
		case "--config":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--config requires a file path")
			}
			i++
			opts.ConfigPath = args[i]
		case "--output":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--output requires a file path")
			}
			i++
			opts.OutputPath = args[i]
		default:
			return opts, fmt.Errorf("unknown option: %s", arg)
		}
	}
	if opts.Quiet && opts.Verbose {
		return opts, fmt.Errorf("--quiet and --verbose cannot be used together")
	}
	if opts.Email && opts.NoEmail {
		return opts, fmt.Errorf("--email and --no-email cannot be used together")
	}
	return opts, nil
}

func isHelpArg(arg string) bool {
	return arg == "--help" || arg == "-h" || arg == "help"
}

type ErrNotImplemented struct {
	Command string
}

func (e ErrNotImplemented) Error() string {
	return fmt.Sprintf("command not implemented yet: %s", e.Command)
}

func IsNotImplemented(err error) bool {
	var target ErrNotImplemented
	return errors.As(err, &target)
}
