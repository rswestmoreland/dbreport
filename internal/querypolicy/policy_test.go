package querypolicy

import "testing"

func TestValidate(t *testing.T) {
	f, p := Defaults(nil, nil)
	if err := Validate("q", "Title", "SELECT 1", SafetyOptions{BlockedFunctions: f, BlockedPatterns: p}); err != nil {
		t.Fatal(err)
	}
	if err := Validate("q", "Title", "WITH x AS (SELECT 1) SELECT * FROM x", SafetyOptions{BlockedFunctions: f, BlockedPatterns: p}); err != nil {
		t.Fatal(err)
	}
	bad := []string{"INSERT INTO t VALUES (1)", "UPDATE t SET a=1", "DELETE FROM t", "DROP TABLE t", "ALTER TABLE t", "CREATE TABLE t(a int)", "TRUNCATE t", "CALL x()", "SET @a=1", "GRANT ALL ON *.* TO u", "REVOKE ALL ON *.* FROM u", "SELECT * INTO OUTFILE '/tmp/x' FROM t", "SELECT * INTO DUMPFILE '/tmp/x' FROM t", "SELECT LOAD_FILE('/tmp/x')", "SELECT SLEEP(1)", "SELECT BENCHMARK(1,MD5('a'))", "SELECT * FROM t FOR UPDATE", "SELECT * FROM t LOCK IN SHARE MODE", "SELECT 1; SELECT 2"}
	for _, q := range bad {
		if Validate("id", "title", q, SafetyOptions{BlockedFunctions: f, BlockedPatterns: p}) == nil {
			t.Fatalf("expected rejection: %s", q)
		}
	}
}
