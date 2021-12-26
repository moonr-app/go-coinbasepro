package coinbasepro_test

import (
	"context"
	"errors"
	"testing"

	"github.com/moonr-app/go-coinbasepro"
)

func TestGetProducts(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	products, err := client.GetProducts(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range products {
		if coinbasepro.StructHasZeroValues(p) && p.StatusMessage != "" {
			t.Fatal("Zero value")
		}
	}
}

func TestGetBook(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	ctx := context.Background()
	_, err := client.GetBook(ctx, "BTC-USD", 1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.GetBook(ctx, "BTC-USD", 2)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.GetBook(ctx, "BTC-USD", 3)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTicker(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	ctx := context.Background()
	ticker, err := client.GetTicker(ctx, "BTC-USD")
	if err != nil {
		t.Fatal(err)
	}

	if coinbasepro.StructHasZeroValues(ticker) {
		t.Fatal("Zero value")
	}

	ticker, err = client.GetTicker(ctx, "ETH-BTC")
	if err != nil {
		t.Fatal(err)
	}

	if coinbasepro.StructHasZeroValues(ticker) {
		t.Fatal(errors.New("Zero value"))
	}

}

func TestListTrades(t *testing.T) {
	var trades []coinbasepro.Trade
	client := coinbasepro.NewTestClient(t)
	cursor := client.ListTrades("BTC-USD", coinbasepro.ListTradesParams{})

	if err := cursor.NextPage(context.Background(), &trades); err != nil {
		t.Fatal(err)
	}

	for _, a := range trades {
		if coinbasepro.StructHasZeroValues(a) {
			t.Fatal(errors.New("Zero value"))
		}
	}
}

func TestGetHistoricRates(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	params := coinbasepro.GetHistoricRatesParams{
		Granularity: 3600,
	}

	historicRates, err := client.GetHistoricRates(context.Background(), "BTC-USD", params)
	if err != nil {
		t.Fatal(err)
	}

	props := []string{"Time", "Low", "High", "Open", "Close", "Volume"}
	if err := coinbasepro.EnsureProperties(historicRates[0], props); err != nil {
		t.Fatal(err)
	}
}

func TestGetStats(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	stats, err := client.GetStats(context.Background(), "BTC-USD")
	if err != nil {
		t.Fatal(err)
	}

	props := []string{"Low", "Open", "Volume", "Last", "Volume30Day"}
	if err := coinbasepro.EnsureProperties(stats, props); err != nil {
		t.Fatal(err)
	}
}
