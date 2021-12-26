package coinbasepro_test

import (
	"context"
	"testing"

	"github.com/moonr-app/go-coinbasepro/v2"
)

func TestGetCurrencies(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	currencies, err := client.GetCurrencies(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range currencies {
		if coinbasepro.StructHasZeroValues(c) {
			t.Fatal("Zero value")
		}
	}
}
