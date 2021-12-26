package coinbasepro_test

import (
	"context"
	"testing"

	"github.com/moonr-app/go-coinbasepro/v2"
)

func TestGetAccounts(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	accounts, err := client.GetAccounts(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Check for decoding issues
	for _, a := range accounts {
		if coinbasepro.StructHasZeroValues(a) {
			t.Fatal("Zero value")
		}
	}
}

func TestGetAccount(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	ctx := context.Background()
	accounts, err := client.GetAccounts(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, a := range accounts {
		account, err := client.GetAccount(ctx, a.ID)
		if err != nil {
			t.Fatal(err)
		}

		// Check for decoding issues
		if coinbasepro.StructHasZeroValues(account) {
			t.Fatal("Zero value")
		}
	}
}
func TestListAccountLedger(t *testing.T) {
	var ledgers []coinbasepro.LedgerEntry
	client := coinbasepro.NewTestClient(t)
	ctx := context.Background()
	accounts, err := client.GetAccounts(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, a := range accounts {
		cursor := client.ListAccountLedger(a.ID)
		for cursor.HasMore {
			if err := cursor.NextPage(ctx, &ledgers); err != nil {
				t.Fatal(err)
			}

			for _, ledger := range ledgers {
				props := []string{"ID", "CreatedAt", "Amount", "Balance", "Type"}
				if err := coinbasepro.EnsureProperties(ledger, props); err != nil {
					t.Fatal(err)
				}

				if ledger.Type == "match" || ledger.Type == "fee" {
					if err := coinbasepro.Ensure(ledger.Details); err != nil {
						t.Fatal("Details is missing")
					}
				}
			}
		}
	}
}

func TestListHolds(t *testing.T) {
	var holds []coinbasepro.Hold
	client := coinbasepro.NewTestClient(t)
	ctx := context.Background()

	accounts, err := client.GetAccounts(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, a := range accounts {
		cursor := client.ListHolds(a.ID)
		for cursor.HasMore {
			if err := cursor.NextPage(ctx, &holds); err != nil {
				t.Fatal(err)
			}

			for _, h := range holds {
				// Check for decoding issues
				if coinbasepro.StructHasZeroValues(h) {
					t.Fatal("Zero value")
				}
			}
		}
	}
}
