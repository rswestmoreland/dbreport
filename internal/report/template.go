package report

import "time"

func formatTime(value time.Time) string {
	return value.Format("2006-01-02 15:04:05")
}

const reportTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}}</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f5f7fb;
      --card: #ffffff;
      --text: #172033;
      --muted: #667085;
      --border: #d9e0ea;
      --accent: #2f6fed;
      --accent-soft: #eaf1ff;
      --danger-soft: #fff3e8;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      font-size: 14px;
      line-height: 1.5;
    }
    .page {
      max-width: 1100px;
      margin: 0 auto;
      padding: 32px 20px 48px;
    }
    .header {
      background: linear-gradient(135deg, #172033, #253858);
      color: #fff;
      border-radius: 18px;
      padding: 28px;
      box-shadow: 0 14px 32px rgba(15, 23, 42, 0.16);
    }
    .header h1 {
      margin: 0 0 8px;
      font-size: 30px;
      letter-spacing: -0.02em;
    }
    .header p { margin: 0; color: #d6deeb; }
     .section {
      background: var(--card);
      border: 1px solid var(--border);
      border-radius: 16px;
      box-shadow: 0 8px 20px rgba(15, 23, 42, 0.05);
    }
         .section {
      margin-top: 14px;
      overflow: hidden;
    }
     .section-header {
      padding: 14px 18px 8px;
      border-bottom: 1px solid var(--border);
    }
    .section-header h2 {
      margin: 0;
      font-size: 20px;
      letter-spacing: -0.01em;
    }
    .meta {
      margin-top: 4px;
      color: var(--muted);
      font-size: 12px;
    }
    .section-body { padding: 10px 18px 14px; }
    .metric-grid { display: grid; grid-template-columns: repeat(auto-fit,minmax(180px,1fr)); gap: 10px; margin-top: 14px; }
    .metric-tile { background: var(--card); border: 1px solid var(--border); border-radius: 12px; padding: 12px 14px; box-shadow: 0 6px 14px rgba(15,23,42,.04); }
    .metric-value {
      font-size: 28px;
      font-weight: 750;
      letter-spacing: -0.03em;
    }
    .metric-label { color: var(--muted); }
    .table-wrap { overflow-x: auto; }
    table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 12px;
      font-size: 13px;
    }
    th, td {
      border-bottom: 1px solid var(--border);
      padding: 9px 10px;
      text-align: left;
      vertical-align: top;
    }
    th {
      background: #f8fafd;
      color: #344054;
      font-weight: 650;
    }
    tr:last-child td { border-bottom: 0; }
    .empty {
      border: 1px dashed var(--border);
      border-radius: 12px;
      color: var(--muted);
      padding: 18px;
      background: #fbfcfe;
    }
    .notice {
      margin-top: 14px;
      border-radius: 10px;
      padding: 10px 12px;
      background: var(--danger-soft);
      color: #8a4b14;
      font-size: 12px;
    }
    .chart {
      width: 100%;
      height: auto;
      display: block;
      margin: 2px 0 8px;
    }
    .chart-bg { fill: #fbfcfe; stroke: #d9e0ea; }
    .chart-label, .chart-value, .chart-tick, .chart-empty-text {
      fill: #475467;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      font-size: 12px;
    }
    .chart-value { font-weight: 650; }
    .chart-bar-fill { fill: var(--accent); }
    .chart-axis { stroke: #b8c2d1; stroke-width: 1; }
    .chart-line-path {
      fill: none;
      stroke: var(--accent);
      stroke-width: 3;
      stroke-linecap: round;
      stroke-linejoin: round;
    }
    .chart-point { fill: #fff; stroke: var(--accent); stroke-width: 2; }
    .footer {
      margin-top: 20px;
      color: var(--muted);
      font-size: 11px;
      text-align: center;
    }
    @media (max-width: 720px) {
      .header { border-radius: 14px; }
      .metric-value { font-size: 34px; }
    }
    @media print {
      body { background: #fff; }
      .page { max-width: none; padding: 0; }
      .header, .summary-card, .section { box-shadow: none; }
      .section { break-inside: avoid; }
    }
  </style>
</head>
<body>
  <main class="page">
    <header class="header">
      <h1>{{.Title}}</h1>
      <p>Generated {{formatTime .GeneratedAt}}</p>
    </header>

    {{range .Blocks}}
    {{if eq .Type "metric-grid"}}
    <section class="metric-grid">
      {{range .Sections}}
      <div class="metric-tile section" id="{{.ID}}">
        <div class="metric-value">{{.MetricValue}}</div>
        <div class="metric-label">{{.Title}}</div>
      </div>
      {{end}}
    </section>
    {{else}}
    {{$s := index .Sections 0}}
    <section class="section" id="{{$s.ID}}">
      <div class="section-header">
        <h2>{{$s.Title}}</h2>
        <div class="meta">{{$s.Type}} section - {{$s.RowCount}} rows - {{$s.Duration}}</div>
      </div>
      <div class="section-body">
        {{if eq $s.Type "bar"}}
          {{$s.ChartHTML}}
          {{if $s.ShowTable}}{{template "table" $s}}{{end}}
        {{else if eq $s.Type "line"}}
          {{$s.ChartHTML}}
          {{if $s.ShowTable}}{{template "table" $s}}{{end}}
        {{else if eq $s.Type "pie"}}
          {{$s.ChartHTML}}
          {{if $s.ShowTable}}{{template "table" $s}}{{end}}
        {{else}}
          {{template "table" $s}}
        {{end}}
        {{if $s.Truncated}}
          <div class="notice">This result was truncated at the configured row cap.</div>
        {{end}}
      </div>
    </section>
    {{end}}
    {{end}}

    <footer class="footer">
      Generated by {{.Footer.AppName}} {{.Footer.Version}} - Build {{.Footer.Commit}} - {{.Footer.Date}} - {{.FooterDetail}}
    </footer>
  </main>
</body>
</html>

{{define "table"}}
  {{if .Rows}}
    <div class="table-wrap">
      <table>
        <thead>
          <tr>{{range .Columns}}<th>{{.}}</th>{{end}}</tr>
        </thead>
        <tbody>
          {{range .Rows}}
            <tr>{{range .}}<td>{{.}}</td>{{end}}</tr>
          {{end}}
        </tbody>
      </table>
    </div>
  {{else}}
    <div class="empty">No rows returned.</div>
  {{end}}
{{end}}
`
