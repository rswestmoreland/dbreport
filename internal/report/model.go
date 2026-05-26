package report

import (
	"html/template"
	"strconv"
	"time"

	"github.com/rswestmoreland/dbreport/internal/config"
	dbreportdb "github.com/rswestmoreland/dbreport/internal/db"
	"github.com/rswestmoreland/dbreport/internal/version"
)

type Document struct {
	Title        string
	GeneratedAt  time.Time
	Blocks       []SectionBlock
	Summary      Summary
	Footer       Footer
	FooterDetail string
}

type Summary struct {
	QueryCount int
	RowCount   int
	Duration   string
}

type SectionBlock struct {
	Type     string
	Sections []Section
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
	ShowTable   bool
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
	summary := Summary{QueryCount: len(results), RowCount: totalRows, Duration: formatDuration(totalDuration)}
	return Document{
		Title:        cfg.Title,
		GeneratedAt:  time.Now(),
		Blocks:       buildBlocks(sections),
		Summary:      summary,
		Footer:       Footer{AppName: version.AppName, Version: version.Version, Commit: version.Commit, Date: version.Date},
		FooterDetail: "Queries: " + itoa(summary.QueryCount) + " - Rows: " + itoa(summary.RowCount) + " - Runtime: " + summary.Duration,
	}, nil
}

func buildBlocks(sections []Section) []SectionBlock {
	blocks := make([]SectionBlock, 0)
	for i := 0; i < len(sections); {
		if sections[i].Type == "metric" {
			j := i
			for j < len(sections) && sections[j].Type == "metric" {
				j++
			}
			blocks = append(blocks, SectionBlock{Type: "metric-grid", Sections: append([]Section(nil), sections[i:j]...)})
			i = j
			continue
		}
		blocks = append(blocks, SectionBlock{Type: "section", Sections: []Section{sections[i]}})
		i++
	}
	return blocks
}

func itoa(v int) string { return strconv.Itoa(v) }
