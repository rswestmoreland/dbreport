package report

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/rswestmoreland/dbreport/internal/charts"
	dbreportdb "github.com/rswestmoreland/dbreport/internal/db"
)

func newSection(result dbreportdb.Result) (Section, error) {
	section := Section{
		ID:        result.Query.ID,
		Title:     result.Query.Title,
		Type:      result.Query.Type,
		Columns:   append([]string(nil), result.Columns...),
		Rows:      stringRows(result),
		Duration:  formatDuration(result.Duration),
		RowCount:  result.RowCount(),
		Truncated: result.Truncated,
		ShowTable: shouldShowTable(result.Query.ShowTable),
	}

	switch strings.ToLower(result.Query.Type) {
	case "metric":
		section.MetricValue, section.MetricLabel = metricValue(result)
	case "table":
		// Tables use the generic row rendering.
	case "bar":
		chart, err := charts.RenderBar(result, result.Query.LabelColumn, result.Query.ValueColumn)
		if err != nil {
			return Section{}, err
		}
		section.ChartHTML = template.HTML(chart)
	case "line":
		chart, err := charts.RenderLine(result, result.Query.LabelColumn, result.Query.SeriesColumn, result.Query.ValueColumn)
		if err != nil {
			return Section{}, err
		}
		section.ChartHTML = template.HTML(chart)
	case "pie":
		chart, err := charts.RenderPie(result, result.Query.LabelColumn, result.Query.ValueColumn)
		if err != nil {
			return Section{}, err
		}
		section.ChartHTML = template.HTML(chart)
	default:
		return Section{}, fmt.Errorf("query %q has unsupported section type %q", result.Query.ID, result.Query.Type)
	}

	return section, nil
}

func stringRows(result dbreportdb.Result) [][]string {
	rows := make([][]string, 0, len(result.Rows))
	for _, sourceRow := range result.Rows {
		row := make([]string, len(sourceRow))
		for i, cell := range sourceRow {
			if cell.IsNull {
				row[i] = "NULL"
			} else {
				row[i] = cell.Text
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func metricValue(result dbreportdb.Result) (string, string) {
	if len(result.Rows) == 0 || len(result.Rows[0]) == 0 {
		return "No data", result.Query.Title
	}
	cell := result.Rows[0][0]
	if cell.IsNull {
		return "NULL", result.Query.Title
	}
	label := result.Query.Title
	if len(result.Columns) > 0 {
		label = result.Columns[0]
	}
	return cell.Text, label
}

func formatDuration(value time.Duration) string {
	if value < time.Millisecond {
		return strings.ReplaceAll(value.String(), "\u00b5", "u")
	}
	return strings.ReplaceAll(value.Round(time.Millisecond).String(), "\u00b5", "u")
}

func shouldShowTable(showTable *bool) bool {
	return showTable == nil || *showTable
}
