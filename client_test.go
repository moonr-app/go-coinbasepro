package coinbasepro_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/preichenberger/go-coinbasepro/v2"
)

func TestClientErrorsOnNotFound(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	_, err := client.Request(context.Background(), http.MethodGet, "/fake", nil, nil)
	if err == nil || err != coinbasepro.ErrNotFound {
		t.Fatal("should have thrown 404 error")
	}
}
