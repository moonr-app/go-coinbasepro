package coinbasepro_test

import (
	"context"
	"testing"
	"time"

	"github.com/preichenberger/go-coinbasepro/v2"
)

func TestCreateReportAndStatus(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	ctx := context.Background()

	newReport := coinbasepro.Report{
		Type:      "fills",
		ProductID: "ALL",
		StartDate: time.Now().Add(-24 * 4 * time.Hour),
		EndDate:   time.Now().Add(-24 * 2 * time.Hour),
	}

	report, err := client.CreateReport(ctx, newReport)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetReportStatus(ctx, report.ID)
	if err != nil {
		t.Fatal(err)
	}
}
