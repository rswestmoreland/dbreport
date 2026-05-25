package db

import (
	"fmt"
	"time"

	"github.com/rswestmoreland/dbreport/internal/config"
)

type Cell struct {
	Raw    any
	Text   string
	IsNull bool
}

type Result struct {
	Query     config.QueryConfig
	Columns   []string
	Rows      [][]Cell
	Duration  time.Duration
	Truncated bool
}

func (r Result) RowCount() int {
	return len(r.Rows)
}

func ConvertCell(value any) Cell {
	switch v := value.(type) {
	case nil:
		return Cell{Raw: nil, Text: "", IsNull: true}
	case []byte:
		return Cell{Raw: v, Text: string(v)}
	case time.Time:
		return Cell{Raw: v, Text: formatTime(v)}
	default:
		return Cell{Raw: v, Text: fmt.Sprint(v)}
	}
}

func formatTime(value time.Time) string {
	return value.Format("2006-01-02 15:04:05")
}
