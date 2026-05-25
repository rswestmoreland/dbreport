package charts

import (
	"fmt"
	"html"
	"math"
	"strconv"
	"strings"

	dbreportdb "github.com/rswestmoreland/dbreport/internal/db"
)

const (
	chartWidth  = 720
	barRowStep  = 36
	barLeftPad  = 150
	barRightPad = 36
	lineHeight  = 260
	lineLeftPad = 54
	lineTopPad  = 24
)

func RenderBar(result dbreportdb.Result, labelColumn string, valueColumn string) (string, error) {
	labelIndex, err := columnIndex(result.Columns, labelColumn)
	if err != nil {
		return "", err
	}
	valueIndex, err := columnIndex(result.Columns, valueColumn)
	if err != nil {
		return "", err
	}

	points, err := chartPoints(result, labelIndex, valueIndex)
	if err != nil {
		return "", err
	}
	if len(points) == 0 {
		return emptySVG("No rows returned")
	}

	maxValue := 0.0
	for _, point := range points {
		if point.Value > maxValue {
			maxValue = point.Value
		}
	}
	if maxValue <= 0 {
		maxValue = 1
	}

	height := 48 + len(points)*barRowStep
	barMaxWidth := chartWidth - barLeftPad - barRightPad
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<svg class="chart chart-bar" viewBox="0 0 %d %d" role="img" aria-label="%s">`, chartWidth, height, html.EscapeString(result.Query.Title)))
	b.WriteString(`<rect class="chart-bg" x="0" y="0" width="100%" height="100%" rx="10"></rect>`)
	for i, point := range points {
		y := 32 + i*barRowStep
		width := int(math.Round((point.Value / maxValue) * float64(barMaxWidth)))
		if point.Value > 0 && width < 1 {
			width = 1
		}
		b.WriteString(fmt.Sprintf(`<text class="chart-label" x="16" y="%d">%s</text>`, y+18, html.EscapeString(trimLabel(point.Label, 22))))
		b.WriteString(fmt.Sprintf(`<rect class="chart-bar-fill" x="%d" y="%d" width="%d" height="20" rx="4"></rect>`, barLeftPad, y, width))
		b.WriteString(fmt.Sprintf(`<text class="chart-value" x="%d" y="%d">%s</text>`, barLeftPad+width+8, y+16, html.EscapeString(formatNumber(point.Value))))
	}
	b.WriteString(`</svg>`)
	return b.String(), nil
}

func RenderLine(result dbreportdb.Result, labelColumn string, valueColumn string) (string, error) {
	labelIndex, err := columnIndex(result.Columns, labelColumn)
	if err != nil {
		return "", err
	}
	valueIndex, err := columnIndex(result.Columns, valueColumn)
	if err != nil {
		return "", err
	}

	points, err := chartPoints(result, labelIndex, valueIndex)
	if err != nil {
		return "", err
	}
	if len(points) == 0 {
		return emptySVG("No rows returned")
	}

	minValue := points[0].Value
	maxValue := points[0].Value
	for _, point := range points {
		if point.Value < minValue {
			minValue = point.Value
		}
		if point.Value > maxValue {
			maxValue = point.Value
		}
	}
	if minValue == maxValue {
		minValue = 0
		if maxValue == 0 {
			maxValue = 1
		}
	}

	plotWidth := chartWidth - lineLeftPad - 32
	plotHeight := lineHeight - lineTopPad - 52
	var polyline strings.Builder
	var markers strings.Builder
	for i, point := range points {
		x := lineLeftPad
		if len(points) > 1 {
			x = lineLeftPad + int(math.Round(float64(i)*float64(plotWidth)/float64(len(points)-1)))
		}
		ratio := (point.Value - minValue) / (maxValue - minValue)
		y := lineTopPad + plotHeight - int(math.Round(ratio*float64(plotHeight)))
		polyline.WriteString(fmt.Sprintf("%d,%d ", x, y))
		markers.WriteString(fmt.Sprintf(`<circle class="chart-point" cx="%d" cy="%d" r="4"></circle>`, x, y))
	}

	last := points[len(points)-1]
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<svg class="chart chart-line" viewBox="0 0 %d %d" role="img" aria-label="%s">`, chartWidth, lineHeight, html.EscapeString(result.Query.Title)))
	b.WriteString(`<rect class="chart-bg" x="0" y="0" width="100%" height="100%" rx="10"></rect>`)
	b.WriteString(fmt.Sprintf(`<line class="chart-axis" x1="%d" y1="%d" x2="%d" y2="%d"></line>`, lineLeftPad, lineTopPad+plotHeight, lineLeftPad+plotWidth, lineTopPad+plotHeight))
	b.WriteString(fmt.Sprintf(`<line class="chart-axis" x1="%d" y1="%d" x2="%d" y2="%d"></line>`, lineLeftPad, lineTopPad, lineLeftPad, lineTopPad+plotHeight))
	b.WriteString(fmt.Sprintf(`<text class="chart-tick" x="12" y="%d">%s</text>`, lineTopPad+12, html.EscapeString(formatNumber(maxValue))))
	b.WriteString(fmt.Sprintf(`<text class="chart-tick" x="12" y="%d">%s</text>`, lineTopPad+plotHeight, html.EscapeString(formatNumber(minValue))))
	b.WriteString(fmt.Sprintf(`<polyline class="chart-line-path" points="%s"></polyline>`, strings.TrimSpace(polyline.String())))
	b.WriteString(markers.String())
	b.WriteString(fmt.Sprintf(`<text class="chart-label" x="%d" y="%d">Latest: %s = %s</text>`, lineLeftPad, lineHeight-18, html.EscapeString(trimLabel(last.Label, 28)), html.EscapeString(formatNumber(last.Value))))
	b.WriteString(`</svg>`)
	return b.String(), nil
}

func emptySVG(message string) (string, error) {
	return fmt.Sprintf(`<svg class="chart chart-empty" viewBox="0 0 %d 120" role="img" aria-label="%s"><rect class="chart-bg" x="0" y="0" width="100%%" height="100%%" rx="10"></rect><text class="chart-empty-text" x="24" y="66">%s</text></svg>`, chartWidth, html.EscapeString(message), html.EscapeString(message)), nil
}

type point struct {
	Label string
	Value float64
}

func chartPoints(result dbreportdb.Result, labelIndex int, valueIndex int) ([]point, error) {
	points := make([]point, 0, len(result.Rows))
	for rowIndex, row := range result.Rows {
		if labelIndex >= len(row) || valueIndex >= len(row) {
			return nil, fmt.Errorf("query %q row %d does not contain required chart columns", result.Query.ID, rowIndex)
		}
		value, err := numericValue(row[valueIndex])
		if err != nil {
			return nil, fmt.Errorf("query %q row %d column %q is not numeric: %w", result.Query.ID, rowIndex, result.Columns[valueIndex], err)
		}
		points = append(points, point{Label: row[labelIndex].Text, Value: value})
	}
	return points, nil
}

func columnIndex(columns []string, name string) (int, error) {
	for i, column := range columns {
		if strings.EqualFold(column, name) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("required column %q was not returned by query", name)
}

func numericValue(cell dbreportdb.Cell) (float64, error) {
	if cell.IsNull {
		return 0, fmt.Errorf("value is NULL")
	}
	switch value := cell.Raw.(type) {
	case int64:
		return float64(value), nil
	case int:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case uint64:
		return float64(value), nil
	case float64:
		return value, nil
	case float32:
		return float64(value), nil
	case []byte:
		return strconv.ParseFloat(strings.TrimSpace(string(value)), 64)
	case string:
		return strconv.ParseFloat(strings.TrimSpace(value), 64)
	default:
		return strconv.ParseFloat(strings.TrimSpace(cell.Text), 64)
	}
}

func trimLabel(value string, maxRunes int) string {
	runes := []rune(value)
	if len(runes) <= maxRunes {
		return value
	}
	if maxRunes <= 1 {
		return "..."
	}
	return string(runes[:maxRunes-1]) + "..."
}

func formatNumber(value float64) string {
	if math.Abs(value-math.Round(value)) < 0.0000001 {
		return strconv.FormatInt(int64(math.Round(value)), 10)
	}
	return strconv.FormatFloat(value, 'f', 2, 64)
}
