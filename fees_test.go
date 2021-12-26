package coinbasepro_test

import (
	"context"
	"testing"

	"github.com/moonr-app/go-coinbasepro"
)

func TestGetFees(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	_, err := client.GetFees(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
