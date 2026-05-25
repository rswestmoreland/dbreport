package report

import (
	"bytes"
	"fmt"
	"html/template"
)

func RenderHTML(doc Document) ([]byte, error) {
	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"formatTime": formatTime,
	}).Parse(reportTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse report template: %w", err)
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, doc); err != nil {
		return nil, fmt.Errorf("render report template: %w", err)
	}
	return out.Bytes(), nil
}
