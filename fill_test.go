package coinbasepro_test

import (
	"context"
	"testing"

	"github.com/moonr-app/go-coinbasepro"
)

func TestListFills(t *testing.T) {
	var fills []coinbasepro.Fill
	client := coinbasepro.NewTestClient(t)
	params := coinbasepro.ListFillsParams{
		ProductID: "BTC-USD",
	}
	cursor := client.ListFills(params)
	for cursor.HasMore {
		if err := cursor.NextPage(context.Background(), &fills); err != nil {
			t.Fatal(err)
		}

		for _, f := range fills {
			if coinbasepro.StructHasZeroValues(f) {
				t.Fatal("Zero value")
			}
		}
	}
}
