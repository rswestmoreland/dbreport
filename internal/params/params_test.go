package params

import "testing"

func TestBindBasic(t *testing.T) {
	b, err := Bind("select * from t where a=:a and b=:b", map[string]any{"a": 1, "b": "x"})
	if err != nil {
		t.Fatal(err)
	}
	if b.SQL != "select * from t where a=? and b=?" {
		t.Fatal(b.SQL)
	}
	if len(b.Args) != 2 {
		t.Fatal(len(b.Args))
	}
}
func TestIgnoreQuotedAndComment(t *testing.T) {
	_, err := Bind("select ':x', \" :y\", `:z` -- :n\n from t where a=:a", map[string]any{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
}
func TestMissing(t *testing.T) {
	_, err := Bind("select * from t where a=:a", map[string]any{})
	if err == nil {
		t.Fatal("expected err")
	}
}
func TestStructuralRejected(t *testing.T) {
	bad := []string{"select * from :t", "select :c from t", "select * from t order by :c", "select * from t group by :c", "select * from t join :x j on 1=1"}
	for _, q := range bad {
		if _, err := Bind(q, map[string]any{"t": "x", "c": "x", "x": "x"}); err == nil {
			t.Fatalf("expected err: %s", q)
		}
	}
}
func TestParseCLI(t *testing.T) {
	_, err := ParseCLI([]string{"a=1", "a=2"})
	if err == nil {
		t.Fatal("dup expected")
	}
	m, err := ParseCLI([]string{"a=", "_b=ok"})
	if err != nil || m["a"] != "" {
		t.Fatal(err, m)
	}
}
