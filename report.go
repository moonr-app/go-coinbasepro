package coinbasepro

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ReportParams struct {
	StartDate time.Time
	EndDate   time.Time
}

type CreateReportParams struct {
	Start time.Time
	End   time.Time
}

type Report struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	// ProductID is required for fills type reports.
	// Use 'ALL' to get all products.
	ProductID string `json:"product_id"`
	// AccountID is required for account type reports.
	AccountID   string       `json:"account_id"`
	CreatedAt   Time         `json:"created_at,string"`
	CompletedAt Time         `json:"completed_at,string,"`
	ExpiresAt   Time         `json:"expires_at,string"`
	FileURL     string       `json:"file_url"`
	Params      ReportParams `json:"params"`
	StartDate   time.Time
	EndDate     time.Time
}

func (c *client) CreateReport(ctx context.Context, newReport Report) (Report, error) {
	var savedReport Report

	url := fmt.Sprintf("/reports")
	_, err := c.Request(ctx, http.MethodPost, url, newReport, &savedReport)

	return savedReport, err
}

func (c *client) GetReportStatus(ctx context.Context, id string) (Report, error) {
	report := Report{}

	url := fmt.Sprintf("/reports/%s", id)
	_, err := c.Request(ctx, http.MethodGet, url, nil, &report)

	return report, err
}
