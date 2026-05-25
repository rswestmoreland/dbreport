package report

import (
	"html/template"
	"time"

	"github.com/rswestmoreland/dbreport/internal/config"
	dbreportdb "github.com/rswestmoreland/dbreport/internal/db"
	"github.com/rswestmoreland/dbreport/internal/version"
)

type Document struct {
	Title       string
	GeneratedAt time.Time
	Sections    []Section
	Summary     Summary
	Footer      Footer
}

type Summary struct {
	QueryCount int
	RowCount   int
	Duration   string
}

type Section struct {
	ID          string
	Title       string
	Type        string
	Columns     []string
	Rows        [][]string
	MetricValue string
	MetricLabel string
	ChartHTML   template.HTML
	Duration    string
	RowCount    int
	Truncated   bool
}

type Footer struct {
	AppName string
	Version string
	Commit  string
	Date    string
}

func NewDocument(cfg config.Config, results []dbreportdb.Result) (Document, error) {
	sections := make([]Section, 0, len(results))
	totalRows := 0
	totalDuration := time.Duration(0)
	for _, result := range results {
		section, err := newSection(result)
		if err != nil {
			return Document{}, err
		}
		sections = append(sections, section)
		totalRows += result.RowCount()
		totalDuration += result.Duration
	}

	return Document{
		Title:       cfg.Title,
		GeneratedAt: time.Now(),
		Sections:    sections,
		Summary: Summary{
			QueryCount: len(results),
			RowCount:   totalRows,
			Duration:   formatDuration(totalDuration),
		},
		Footer: Footer{
			AppName: version.AppName,
			Version: version.Version,
			Commit:  version.Commit,
			Date:    version.Date,
		},
	}, nil
}
