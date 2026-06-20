package repository

import (
	"context"
	"time"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func (db *Database) GetMetricsDashboard(ctx context.Context, month string, startDate, endDate time.Time, qr *bool) (*domain.MetricsDashboard, error) {
	params := metricsParams(startDate, endDate, qr)

	totals, err := db.Q.GetMetricsTotals(ctx, sqlc.GetMetricsTotalsParams(params))
	if err != nil {
		return nil, err
	}

	bookRows, err := db.Q.GetBookMetricsSummary(ctx, sqlc.GetBookMetricsSummaryParams(params))
	if err != nil {
		return nil, err
	}

	machineRows, err := db.Q.GetMachineMetricsSummary(ctx, sqlc.GetMachineMetricsSummaryParams(params))
	if err != nil {
		return nil, err
	}

	return &domain.MetricsDashboard{
		Month: month,
		Qr:    qrLabel(qr),
		Totals: domain.MetricsTotals{
			BookClicks:     totals.BookClicks,
			MachineViews:   totals.MachineViews,
			UniqueSessions: totals.UniqueSessions,
		},
		Books:    mapBookMetricSummaries(bookRows),
		Machines: mapMachineMetricSummaries(machineRows),
	}, nil
}

type metricsQueryParams struct {
	StartDate pgtype.Date
	EndDate   pgtype.Date
	Qr        pgtype.Bool
}

func metricsParams(startDate, endDate time.Time, qr *bool) metricsQueryParams {
	params := metricsQueryParams{
		StartDate: pgtype.Date{Time: startDate, Valid: true},
		EndDate:   pgtype.Date{Time: endDate, Valid: true},
	}
	if qr != nil {
		params.Qr = pgtype.Bool{Bool: *qr, Valid: true}
	}
	return params
}

func qrLabel(qr *bool) string {
	if qr == nil {
		return "all"
	}
	if *qr {
		return "true"
	}
	return "false"
}

func mapBookMetricSummaries(rows []sqlc.GetBookMetricsSummaryRow) []domain.BookMetricsSummary {
	out := make([]domain.BookMetricsSummary, len(rows))
	for i, row := range rows {
		out[i] = domain.BookMetricsSummary{
			BookID:         row.BookID,
			Title:          row.Title,
			Author:         row.Author,
			Clicks:         row.Clicks,
			UniqueSessions: row.UniqueSessions,
		}
	}
	return out
}

func mapMachineMetricSummaries(rows []sqlc.GetMachineMetricsSummaryRow) []domain.MachineMetricsSummary {
	out := make([]domain.MachineMetricsSummary, len(rows))
	for i, row := range rows {
		out[i] = domain.MachineMetricsSummary{
			MachineID:      row.MachineID,
			Location:       row.Location,
			Views:          row.Views,
			BookClicks:     row.BookClicks,
			UniqueSessions: row.UniqueSessions,
		}
	}
	return out
}
