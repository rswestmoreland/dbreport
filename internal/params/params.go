package params

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
)

var nameRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type BindResult struct {
	SQL           string
	Args          []any
	ReferencedSet map[string]struct{}
}

func ValidateName(name string) error {
	if !nameRE.MatchString(name) {
		return fmt.Errorf("invalid parameter name %q: must match [A-Za-z_][A-Za-z0-9_]*", name)
	}
	return nil
}
func ParseCLI(items []string) (map[string]any, error) {
	out := map[string]any{}
	for _, it := range items {
		eq := strings.IndexByte(it, '=')
		if eq < 0 {
			return nil, fmt.Errorf("invalid --param %q: expected name=value", it)
		}
		n, v := it[:eq], it[eq+1:]
		if err := ValidateName(n); err != nil {
			return nil, err
		}
		if _, ok := out[n]; ok {
			return nil, fmt.Errorf("duplicate --param name: %q", n)
		}
		out[n] = v
	}
	return out, nil
}

func LoadYAMLFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read params file %q: %w", path, err)
	}
	var dec any
	if err := yaml.Unmarshal(data, &dec); err != nil {
		return nil, fmt.Errorf("parse params file %q: %w", path, err)
	}
	m, ok := dec.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("params file %q must contain a top-level mapping", path)
	}
	out := map[string]any{}
	for k, v := range m {
		if err := ValidateName(k); err != nil {
			return nil, err
		}
		switch v.(type) {
		case nil, string, int, int64, uint64, float64, bool:
			out[k] = v
		default:
			return nil, fmt.Errorf("params file %q value for %q must be a scalar", path, k)
		}
	}
	return out, nil
}
func Merge(file, cli map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range file {
		out[k] = v
	}
	for k, v := range cli {
		out[k] = v
	}
	return out
}

func Bind(sqlText string, values map[string]any) (BindResult, error) {
	var b strings.Builder
	args := []any{}
	used := map[string]struct{}{}
	inSQ, inDQ, inBT, inLine, inBlock := false, false, false, false, false
	for i := 0; i < len(sqlText); i++ {
		c := sqlText[i]
		if inLine {
			b.WriteByte(c)
			if c == '\n' {
				inLine = false
			}
			continue
		}
		if inBlock {
			b.WriteByte(c)
			if c == '*' && i+1 < len(sqlText) && sqlText[i+1] == '/' {
				i++
				b.WriteByte('/')
				inBlock = false
			}
			continue
		}
		if inSQ {
			b.WriteByte(c)
			if c == '\'' {
				if i+1 < len(sqlText) && sqlText[i+1] == '\'' {
					i++
					b.WriteByte('\'')
				} else {
					inSQ = false
				}
			}
			continue
		}
		if inDQ {
			b.WriteByte(c)
			if c == '"' {
				inDQ = false
			}
			continue
		}
		if inBT {
			b.WriteByte(c)
			if c == '`' {
				inBT = false
			}
			continue
		}
		if c == '-' && i+1 < len(sqlText) && sqlText[i+1] == '-' {
			inLine = true
			b.WriteString("--")
			i++
			continue
		}
		if c == '/' && i+1 < len(sqlText) && sqlText[i+1] == '*' {
			inBlock = true
			b.WriteString("/*")
			i++
			continue
		}
		if c == '\'' {
			inSQ = true
			b.WriteByte(c)
			continue
		}
		if c == '"' {
			inDQ = true
			b.WriteByte(c)
			continue
		}
		if c == '`' {
			inBT = true
			b.WriteByte(c)
			continue
		}
		if c == ':' && i+1 < len(sqlText) && isNameStart(sqlText[i+1]) {
			j := i + 2
			for j < len(sqlText) && isNamePart(sqlText[j]) {
				j++
			}
			n := sqlText[i+1 : j]
			if err := validateValueContext(sqlText, i, j); err != nil {
				return BindResult{}, fmt.Errorf("parameter %q used in SQL structure: %w", n, err)
			}
			v, ok := values[n]
			if !ok {
				return BindResult{}, fmt.Errorf("parameter %q is required, but no value was provided", n)
			}
			b.WriteByte('?')
			args = append(args, v)
			used[n] = struct{}{}
			i = j - 1
			continue
		}
		b.WriteByte(c)
	}
	return BindResult{SQL: b.String(), Args: args, ReferencedSet: used}, nil
}

func ValidateUnused(all map[string]any, used map[string]struct{}) error {
	var u []string
	for k := range all {
		if _, ok := used[k]; !ok {
			u = append(u, k)
		}
	}
	if len(u) == 0 {
		return nil
	}
	sort.Strings(u)
	return fmt.Errorf("unused parameter(s): %s", strings.Join(u, ", "))
}
func isNameStart(c byte) bool { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' }
func isNamePart(c byte) bool  { return isNameStart(c) || (c >= '0' && c <= '9') }
func validateValueContext(s string, start, end int) error {
	up := strings.ToUpper(s)
	pre := strings.ToUpper(prevWord(s, start-1))
	if pre == "FROM" || pre == "JOIN" || pre == "TABLE" || pre == "DATABASE" {
		return fmt.Errorf("identifier position")
	}
	if strings.Contains(up[max(0, start-20):start], "ORDER BY") {
		return fmt.Errorf("ORDER BY expression")
	}
	if strings.Contains(up[max(0, start-20):start], "GROUP BY") {
		return fmt.Errorf("GROUP BY expression")
	}
	if prevNonSpace(s, start-1) == "." || nextNonSpace(s, end) == "." {
		return fmt.Errorf("dynamic identifier")
	}
	if inSelectList(up, start) {
		return fmt.Errorf("SELECT list")
	}
	return nil
}
func prevWord(s string, i int) string {
	for i >= 0 && s[i] <= ' ' {
		i--
	}
	j := i
	for j >= 0 && ((s[j] >= 'A' && s[j] <= 'Z') || (s[j] >= 'a' && s[j] <= 'z') || s[j] == '_') {
		j--
	}
	return s[j+1 : i+1]
}
func prevNonSpace(s string, i int) string {
	for ; i >= 0; i-- {
		if s[i] > ' ' {
			return string(s[i])
		}
	}
	return ""
}
func nextNonSpace(s string, i int) string {
	for ; i < len(s); i++ {
		if s[i] > ' ' {
			return string(s[i])
		}
	}
	return ""
}
func inSelectList(up string, pos int) bool {
	sel := strings.LastIndex(up[:pos], "SELECT")
	if sel < 0 {
		return false
	}
	return strings.Index(up[sel:pos], "FROM") < 0
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
