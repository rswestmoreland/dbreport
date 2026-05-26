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

	rendered, err := RenderLine(result, "day", "", "count")
	if err != nil {
		t.Fatalf("RenderLine failed: %v", err)
	}
	for _, want := range []string{"chart-line", "polyline", "2026-05-24", "2026-05-25", ">0<", "text-anchor=\"start\"", "text-anchor=\"end\""} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("expected line chart to contain %q, got %s", want, rendered)
		}
	}
}

func TestRenderLineMultiSeriesLegend(t *testing.T) {
	result := dbreportdb.Result{
		Query:   config.QueryConfig{ID: "daily", Title: "Daily", Type: "line"},
		Columns: []string{"day", "result", "count"},
		Rows: [][]dbreportdb.Cell{
			{dbreportdb.ConvertCell("2026-05-24"), dbreportdb.ConvertCell("success"), dbreportdb.ConvertCell(int64(8))},
			{dbreportdb.ConvertCell("2026-05-24"), dbreportdb.ConvertCell("failure"), dbreportdb.ConvertCell(int64(2))},
			{dbreportdb.ConvertCell("2026-05-25"), dbreportdb.ConvertCell("success"), dbreportdb.ConvertCell(int64(9))},
			{dbreportdb.ConvertCell("2026-05-25"), dbreportdb.ConvertCell("failure"), dbreportdb.ConvertCell(int64(4))},
		},
	}
	rendered, err := RenderLine(result, "day", "result", "count")
	if err != nil {
		t.Fatalf("RenderLine failed: %v", err)
	}
	for _, want := range []string{"success", "failure", "stroke-linecap:round"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("expected legend marker %q in %s", want, rendered)
		}
	}
	if strings.Contains(rendered, "...") {
		t.Fatalf("labels should not be truncated: %s", rendered)
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

	_, err := RenderLine(result, "day", "", "count")
	if err == nil {
		t.Fatal("expected non-numeric error")
	}
	if !strings.Contains(err.Error(), "is not numeric") {
		t.Fatalf("expected non-numeric error, got %v", err)
	}
}

func TestRenderPieChart(t *testing.T) {
	result := dbreportdb.Result{
		Query:   config.QueryConfig{ID: "share", Title: "Share", Type: "pie"},
		Columns: []string{"result", "count"},
		Rows: [][]dbreportdb.Cell{
			{dbreportdb.ConvertCell("success"), dbreportdb.ConvertCell(int64(3))},
			{dbreportdb.ConvertCell("failure"), dbreportdb.ConvertCell(int64(2))},
		},
	}
	rendered, err := RenderPie(result, "result", "count")
	if err != nil || !strings.Contains(rendered, "chart-pie") || !strings.Contains(rendered, "success: 3") || !strings.Contains(rendered, "stroke-width=\"2\"") {
		t.Fatalf("unexpected pie output err=%v output=%s", err, rendered)
	}
	if !strings.Contains(rendered, "L 190.00 24.00") {
		t.Fatalf("expected first slice to start at top boundary, output=%s", rendered)
	}
	if strings.Index(rendered, "success: 3") > strings.Index(rendered, "failure: 2") {
		t.Fatalf("expected legend to preserve query row order, output=%s", rendered)
	}
}

func TestRenderPieRejectsNegativeValues(t *testing.T) {
	result := dbreportdb.Result{
		Query:   config.QueryConfig{ID: "share", Title: "Share", Type: "pie"},
		Columns: []string{"result", "count"},
		Rows: [][]dbreportdb.Cell{
			{dbreportdb.ConvertCell("success"), dbreportdb.ConvertCell(int64(3))},
			{dbreportdb.ConvertCell("failure"), dbreportdb.ConvertCell(int64(-2))},
		},
	}

	_, err := RenderPie(result, "result", "count")
	if err == nil {
		t.Fatal("expected negative pie value error")
	}
	if !strings.Contains(err.Error(), "negative pie values") {
		t.Fatalf("unexpected error: %v", err)
	}
}
