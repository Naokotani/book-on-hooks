package domain

type MetricsTotals struct {
	BookClicks     int64 `json:"book_clicks"`
	MachineViews   int64 `json:"machine_views"`
	UniqueSessions int64 `json:"unique_sessions"`
}

type BookMetricsSummary struct {
	BookID         int64  `json:"book_id"`
	Title          string `json:"title"`
	Author         string `json:"author"`
	Clicks         int64  `json:"clicks"`
	UniqueSessions int64  `json:"unique_sessions"`
}

type MachineMetricsSummary struct {
	MachineID      int64  `json:"machine_id"`
	Location       string `json:"location"`
	Views          int64  `json:"views"`
	BookClicks     int64  `json:"book_clicks"`
	UniqueSessions int64  `json:"unique_sessions"`
}

type MetricsDashboard struct {
	Month    string                  `json:"month"`
	Qr       string                  `json:"qr"`
	Totals   MetricsTotals           `json:"totals"`
	Books    []BookMetricsSummary    `json:"books"`
	Machines []MachineMetricsSummary `json:"machines"`
}
