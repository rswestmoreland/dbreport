package querypolicy

import (
	"fmt"
	"regexp"
	"strings"
)

var commentBlock = regexp.MustCompile(`(?s)/\*.*?\*/`)
var commentLine = regexp.MustCompile(`(?m)--[^\n]*$`)
var ws = regexp.MustCompile(`\s+`)

type SafetyOptions struct {
	ActiveDatabase                    string
	AllowedDatabases, AllowedTables   []string
	BlockedFunctions, BlockedPatterns []string
}

var defaultBlockedFunctions = []string{"sleep", "benchmark", "load_file"}
var defaultBlockedPatterns = []string{`(?i)password`, `(?i)passwd`, `(?i)token`, `(?i)api[_-]?key`, `(?i)secret`, `(?i)private[_-]?key`, `(?i)session`, `(?i)cookie`}

func Defaults(functions, patterns []string) (outFns, outPats []string) {
	outFns = defaultBlockedFunctions
	outPats = defaultBlockedPatterns
	if len(functions) > 0 {
		outFns = functions
	}
	if len(patterns) > 0 {
		outPats = patterns
	}
	return
}

func Validate(queryID, title, sql string, opts SafetyOptions) error {
	norm := normalize(sql)
	if norm == "" {
		return scopedErr(queryID, title, "query is empty")
	}
	if strings.Count(norm, ";") > 0 && !strings.HasSuffix(norm, ";") {
		return scopedErr(queryID, title, "query policy rejected SQL: multiple statements are not allowed")
	}
	norm = strings.TrimSuffix(norm, ";")
	low := strings.ToLower(norm)
	if !(strings.HasPrefix(low, "select ") || strings.HasPrefix(low, "with ")) {
		return scopedErr(queryID, title, "query policy rejected SQL: only SELECT or WITH ... SELECT are allowed")
	}
	banned := []string{" insert ", " update ", " delete ", " drop ", " alter ", " create ", " truncate ", " replace ", " call ", " do ", " set ", " grant ", " revoke ", " into outfile", " into dumpfile", " load_file(", " sleep(", " benchmark(", " for update", " lock in share mode"}
	padded := " " + low + " "
	for _, b := range banned {
		if strings.Contains(padded, b) {
			return scopedErr(queryID, title, fmt.Sprintf("query policy rejected SQL: contains blocked token %q", strings.TrimSpace(b)))
		}
	}
	for _, fn := range opts.BlockedFunctions {
		if strings.Contains(low, strings.ToLower(fn)+"(") {
			return scopedErr(queryID, title, fmt.Sprintf("query policy rejected SQL: blocked function %q", fn))
		}
	}
	for _, p := range opts.BlockedPatterns {
		if re, err := regexp.Compile(p); err == nil && re.MatchString(norm) {
			return scopedErr(queryID, title, fmt.Sprintf("query policy rejected SQL: matched blocked pattern %q", p))
		}
	}
	if len(opts.AllowedTables) > 0 {
		if err := validateTables(queryID, title, low, opts.AllowedTables); err != nil {
			return err
		}
	}
	if len(opts.AllowedDatabases) > 0 {
		if err := validateDatabases(queryID, title, low, opts.ActiveDatabase, opts.AllowedDatabases); err != nil {
			return err
		}
	}
	return nil
}
func normalize(s string) string {
	s = commentBlock.ReplaceAllString(s, " ")
	s = commentLine.ReplaceAllString(s, " ")
	s = ws.ReplaceAllString(strings.TrimSpace(s), " ")
	return s
}
func scopedErr(id, title, msg string) error {
	if title != "" {
		return fmt.Errorf("query %q (%s): %s", id, title, msg)
	}
	return fmt.Errorf("query %q: %s", id, msg)
}
func validateTables(id, title, low string, allowed []string) error {
	if strings.Contains(low, " from (") {
		return scopedErr(id, title, "query policy rejected SQL: query too complex to validate allowed_tables")
	}
	allow := map[string]struct{}{}
	for _, t := range allowed {
		allow[strings.ToLower(t)] = struct{}{}
	}
	re := regexp.MustCompile(`\b(?:from|join)\s+([a-zA-Z0-9_\.]+)`)
	m := re.FindAllStringSubmatch(low, -1)
	if len(m) == 0 {
		return scopedErr(id, title, "query policy rejected SQL: could not validate allowed_tables")
	}
	for _, g := range m {
		tbl := strings.Trim(g[1], "`")
		parts := strings.Split(tbl, ".")
		tbl = parts[len(parts)-1]
		if _, ok := allow[tbl]; !ok {
			return scopedErr(id, title, fmt.Sprintf("query policy rejected SQL: table %q is not in safety.allowed_tables", tbl))
		}
	}
	return nil
}
func validateDatabases(id, title, low, active string, allowed []string) error {
	if strings.Contains(low, " from (") || strings.Contains(low, " join (") {
		return scopedErr(id, title, "query policy rejected SQL: query too complex to validate allowed_databases")
	}
	allow := map[string]struct{}{}
	for _, d := range allowed {
		allow[strings.ToLower(strings.TrimSpace(d))] = struct{}{}
	}
	re := regexp.MustCompile(`\b(?:from|join)\s+([a-zA-Z0-9_\.]+)`)
	m := re.FindAllStringSubmatch(low, -1)
	if len(m) == 0 {
		return scopedErr(id, title, "query policy rejected SQL: could not validate allowed_databases")
	}
	active = strings.ToLower(strings.TrimSpace(active))
	for _, g := range m {
		ref := strings.Trim(g[1], "`")
		parts := strings.Split(ref, ".")
		db := active
		if len(parts) > 1 {
			db = strings.ToLower(parts[0])
		}
		if db == "" {
			return scopedErr(id, title, "query policy rejected SQL: could not determine database for allowed_databases enforcement")
		}
		if _, ok := allow[db]; !ok {
			return scopedErr(id, title, fmt.Sprintf("query policy rejected SQL: database %q is not in safety.allowed_databases", db))
		}
	}
	return nil
}
