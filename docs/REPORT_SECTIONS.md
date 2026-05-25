# dbreport Report Sections

`dbreport` turns each configured query into one report section.

Supported section types:

```text
metric
table
bar
line
```

## Metric

Use for a single number or single value.

```yaml
- id: "total_orders"
  title: "Total Orders"
  type: "metric"
  sql: |
    SELECT COUNT(*) AS value
    FROM orders
```

Contract:

- First row is used.
- First column is used unless a later phase adds explicit value selection.
- Empty result renders as an empty/placeholder metric.

## Table

Use for tabular detail.

```yaml
- id: "recent_orders"
  title: "Recent Orders"
  type: "table"
  sql: |
    SELECT id, status, created_at
    FROM orders
    ORDER BY created_at DESC
    LIMIT 25
```

Contract:

- Column headers come from query result columns.
- Values are HTML-escaped.
- Empty results render a clear empty-state message.
- Global row cap applies.

## Bar chart

Use for categorical counts or rankings.

```yaml
- id: "orders_by_status"
  title: "Orders by Status"
  type: "bar"
  label_column: "status"
  value_column: "count"
  sql: |
    SELECT status, COUNT(*) AS count
    FROM orders
    GROUP BY status
    ORDER BY count DESC
```

Contract:

- `label_column` is required.
- `value_column` is required.
- Values must be numeric or convertible to numeric values.
- Inline SVG is rendered.
- A fallback data table is included.

## Line chart

Use for time series or ordered trends.

```yaml
- id: "daily_orders"
  title: "Daily Orders"
  type: "line"
  label_column: "day"
  value_column: "count"
  sql: |
    SELECT DATE(created_at) AS day, COUNT(*) AS count
    FROM orders
    WHERE created_at >= NOW() - INTERVAL 30 DAY
    GROUP BY DATE(created_at)
    ORDER BY day
```

Contract:

- `label_column` is required.
- `value_column` is required.
- Query order is preserved.
- Values must be numeric or convertible to numeric values.
- Inline SVG is rendered.
- A fallback data table is included.

## Query design recommendations

- Keep queries read-only.
- Prefer explicit column aliases.
- Use `ORDER BY` for charts.
- Use `LIMIT` for detail tables.
- Keep expensive aggregation queries indexed.
- Avoid selecting sensitive fields unless necessary.

## Layout behavior

- Query order in YAML controls report order.
- Consecutive `metric` sections are automatically rendered as a compact metric-tile grid.
- Non-metric sections remain full-width cards in the same configured order.
- Metric tile size and placement are not YAML-configurable yet.
