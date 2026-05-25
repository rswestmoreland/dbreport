package charts

import (
	"strings"
	"testing"
	"time"

	"github.com/rswestmoreland/dbreport/internal/config"
	dbreportdb "github.com/rswestmoreland/dbreport/internal/db"
)

func TestRenderBarEscapesLabels(t *testing.T) {
	result := dbreportdb.Result{
		Query:    config.QueryConfig{ID: "status", Title: "Status", Type: "bar"},
		Columns:  []string{"status", "count"},
		Rows:     [][]dbreportdb.Cell{{dbreportdb.ConvertCell("<open>"), dbreportdb.ConvertCell(int64(3))}},
		Duration: time.Millisecond,
	}

	rendered, err := RenderBar(result, "status", "count")
	if err != nil {
		t.Fatalf("RenderBar failed: %v", err)
	}
	if !strings.Contains(rendered, "&lt;open&gt;") {
		t.Fatalf("expected escaped label, got %s", rendered)
	}
	if strings.Contains(rendered, "<open>") {
		t.Fatalf("expected raw label to be escaped, got %s", rendered)
	}
}

func TestRenderLineChart(t *testing.T) {
	result := dbreportdb.Result{
		Query:    config.QueryConfig{ID: "daily", Title: "Daily", Type: "line"},
		Columns:  []string{"day", "count"},
		Rows:     [][]dbreportdb.Cell{{dbreportdb.ConvertCell("2026-05-24"), dbreportdb.ConvertCell(int64(4))}, {dbreportdb.ConvertCell("2026-05-25"), dbreportdb.ConvertCell(int64(8))}},
		Duration: time.Millisecond,
	}

	rendered, err := RenderLine(result, "day", "count")
	if err != nil {
		t.Fatalf("RenderLine failed: %v", err)
	}
	for _, want := range []string{"chart-line", "polyline", "Latest:"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("expected line chart to contain %q, got %s", want, rendered)
		}
	}
}

func TestRenderBarRejectsMissingColumn(t *testing.T) {
	result := dbreportdb.Result{
		Query:   config.QueryConfig{ID: "status", Title: "Status", Type: "bar"},
		Columns: []string{"status", "count"},
		Rows:    [][]dbreportdb.Cell{{dbreportdb.ConvertCell("open"), dbreportdb.ConvertCell(int64(3))}},
	}

	_, err := RenderBar(result, "missing", "count")
	if err == nil {
		t.Fatal("expected missing column error")
	}
	if !strings.Contains(err.Error(), "required column") {
		t.Fatalf("expected required column error, got %v", err)
	}
}

func TestRenderLineRejectsNonNumericValues(t *testing.T) {
	result := dbreportdb.Result{
		Query:   config.QueryConfig{ID: "daily", Title: "Daily", Type: "line"},
		Columns: []string{"day", "count"},
		Rows:    [][]dbreportdb.Cell{{dbreportdb.ConvertCell("2026-05-25"), dbreportdb.ConvertCell("not-a-number")}},
	}

	_, err := RenderLine(result, "day", "count")
	if err == nil {
		t.Fatal("expected non-numeric error")
	}
	if !strings.Contains(err.Error(), "is not numeric") {
		t.Fatalf("expected non-numeric error, got %v", err)
	}
}
