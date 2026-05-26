package db

import (
	"strings"
	"testing"

	"github.com/rswestmoreland/dbreport/internal/config"
)

func TestValidateReturnedColumns(t *testing.T) {
	res := Result{Query: config.QueryConfig{ID: "q", Title: "Query"}, Columns: []string{"email", "password_hash"}}
	if err := ValidateReturnedColumns(res, []string{"password_hash"}, []string{}); err == nil || !strings.Contains(err.Error(), "password_hash") {
		t.Fatalf("expected blocked column error, got %v", err)
	}
	res2 := Result{Query: config.QueryConfig{ID: "q2", Title: "Query 2"}, Columns: []string{"apiToken"}}
	if err := ValidateReturnedColumns(res2, []string{}, []string{`(?i)token`}); err == nil || !strings.Contains(err.Error(), "apiToken") {
		t.Fatalf("expected blocked pattern error, got %v", err)
	}
	res3 := Result{Query: config.QueryConfig{ID: "q3", Title: "Query 3"}, Columns: []string{"email", "last_login"}}
	if err := ValidateReturnedColumns(res3, []string{"password"}, []string{`(?i)secret`}); err != nil {
		t.Fatalf("expected safe columns to pass, got %v", err)
	}
}
