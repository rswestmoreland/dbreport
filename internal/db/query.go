package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rswestmoreland/dbreport/internal/config"
)

func RunAll(ctx context.Context, handle *sql.DB, queries []config.QueryConfig, maxRows int, timeout time.Duration) ([]Result, error) {
	results := make([]Result, 0, len(queries))
	for _, query := range queries {
		result, err := RunQuery(ctx, handle, query, maxRows, timeout)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}
	return results, nil
}

func RunQuery(ctx context.Context, handle *sql.DB, query config.QueryConfig, maxRows int, timeout time.Duration) (Result, error) {
	if maxRows <= 0 {
		return Result{}, fmt.Errorf("max rows per query must be greater than zero")
	}
	if timeout <= 0 {
		return Result{}, fmt.Errorf("query timeout must be greater than zero")
	}

	queryCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	rows, err := handle.QueryContext(queryCtx, query.SQL)
	if err != nil {
		return Result{}, QueryError{QueryID: query.ID, Title: query.Title, Err: err}
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return Result{}, QueryError{QueryID: query.ID, Title: query.Title, Err: err}
	}

	result := Result{
		Query:   query,
		Columns: columns,
		Rows:    make([][]Cell, 0),
	}

	for rows.Next() {
		if len(result.Rows) >= maxRows {
			result.Truncated = true
			break
		}

		values := make([]any, len(columns))
		scanTargets := make([]any, len(columns))
		for i := range values {
			scanTargets[i] = &values[i]
		}

		if err := rows.Scan(scanTargets...); err != nil {
			return result, QueryError{QueryID: query.ID, Title: query.Title, Err: err}
		}

		row := make([]Cell, len(columns))
		for i, value := range values {
			row[i] = ConvertCell(value)
		}
		result.Rows = append(result.Rows, row)
	}

	if err := rows.Err(); err != nil {
		return result, QueryError{QueryID: query.ID, Title: query.Title, Err: err}
	}

	result.Duration = time.Since(start)
	return result, nil
}

type QueryError struct {
	QueryID string
	Title   string
	Err     error
}

func (e QueryError) Error() string {
	if e.Title == "" {
		return fmt.Sprintf("query %q failed: %v", e.QueryID, e.Err)
	}
	return fmt.Sprintf("query %q (%s) failed: %v", e.QueryID, e.Title, e.Err)
}

func (e QueryError) Unwrap() error {
	return e.Err
}
