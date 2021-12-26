package coinbasepro

import (
	"errors"
	"testing"
)

func TestGetCurrencies(t *testing.T) {
	client := NewTestClient(t)
	currencies, err := client.GetCurrencies()
	if err != nil {
		t.Error(err)
	}

	for _, c := range currencies {
		if StructHasZeroValues(c) {
			t.Error(errors.New("Zero value"))
		}
	}
}
