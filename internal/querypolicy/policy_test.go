package querypolicy

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	f, p := Defaults(nil, nil)
	if err := Validate("q", "Title", "SELECT 1", SafetyOptions{BlockedFunctions: f, BlockedPatterns: p}); err != nil {
		t.Fatal(err)
	}
	if err := Validate("q", "Title", "WITH x AS (SELECT 1) SELECT * FROM x", SafetyOptions{BlockedFunctions: f, BlockedPatterns: p}); err != nil {
		t.Fatal(err)
	}
	bad := []string{"INSERT INTO t VALUES (1)", "UPDATE t SET a=1", "DELETE FROM t", "DROP TABLE t", "ALTER TABLE t", "CREATE TABLE t(a int)", "TRUNCATE t", "CALL x()", "SET @a=1", "GRANT ALL ON *.* TO u", "REVOKE ALL ON *.* FROM u", "SELECT * INTO OUTFILE '/tmp/x' FROM t", "SELECT * INTO DUMPFILE '/tmp/x' FROM t", "SELECT LOAD_FILE('/tmp/x')", "SELECT SLEEP(1)", "SELECT BENCHMARK(1,MD5('a'))", "SELECT * FROM t FOR UPDATE", "SELECT * FROM t LOCK IN SHARE MODE", "SELECT 1; SELECT 2", "\nUpDaTe t set a=1"}
	for _, q := range bad {
		if Validate("id", "title", q, SafetyOptions{BlockedFunctions: f, BlockedPatterns: p}) == nil {
			t.Fatalf("expected rejection: %s", q)
		}
	}
}

func TestAllowedDatabases(t *testing.T) {
	opts := SafetyOptions{ActiveDatabase: "dbreport_test", AllowedDatabases: []string{"dbreport_test"}}
	if err := Validate("q1", "t", "SELECT * FROM user_accounts", opts); err != nil {
		t.Fatal(err)
	}
	if err := Validate("q2", "t", "SELECT * FROM dbreport_test.user_accounts", opts); err != nil {
		t.Fatal(err)
	}
	if err := Validate("q3", "t", "SELECT * FROM otherdb.user_accounts", opts); err == nil || !strings.Contains(err.Error(), "allowed_databases") {
		t.Fatalf("expected allowed_databases rejection, got %v", err)
	}
	if err := Validate("q4", "t", "SELECT * FROM (SELECT 1) x", opts); err == nil || !strings.Contains(err.Error(), "too complex") {
		t.Fatalf("expected complexity rejection, got %v", err)
	}
}
