package report

import (
	"strings"
	"testing"
	"time"

	"github.com/rswestmoreland/dbreport/internal/config"
	dbreportdb "github.com/rswestmoreland/dbreport/internal/db"
)

func TestRenderHTMLEscapesTableValues(t *testing.T) {
	cfg := config.Config{Title: "Escape Report"}
	results := []dbreportdb.Result{
		{
			Query:    config.QueryConfig{ID: "rows", Title: "Rows", Type: "table"},
			Columns:  []string{"name"},
			Rows:     [][]dbreportdb.Cell{{dbreportdb.ConvertCell("<script>alert(1)</script>")}},
			Duration: time.Millisecond,
		},
	}

	doc, err := NewDocument(cfg, results)
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	rendered, err := RenderHTML(doc)
	if err != nil {
		t.Fatalf("RenderHTML failed: %v", err)
	}
	text := string(rendered)
	if strings.Contains(text, "<script>alert(1)</script>") {
		t.Fatalf("expected script value to be escaped, got %s", text)
	}
	if !strings.Contains(text, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Fatalf("expected escaped script value, got %s", text)
	}
}

func TestRenderHTMLContainsNoExternalAssetReferences(t *testing.T) {
	cfg := config.Config{Title: "Self Contained Report"}
	results := []dbreportdb.Result{
		{
			Query:    config.QueryConfig{ID: "total", Title: "Total", Type: "metric"},
			Columns:  []string{"value"},
			Rows:     [][]dbreportdb.Cell{{dbreportdb.ConvertCell(int64(1))}},
			Duration: time.Millisecond,
		},
	}

	doc, err := NewDocument(cfg, results)
	if err != nil {
		t.Fatalf("NewDocument failed: %v", err)
	}
	rendered, err := RenderHTML(doc)
	if err != nil {
		t.Fatalf("RenderHTML failed: %v", err)
	}
	lower := strings.ToLower(string(rendered))
	for _, forbidden := range []string{"http://", "https://", "<script", " src=", " href="} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("expected no external reference marker %q in rendered HTML", forbidden)
		}
	}
}
