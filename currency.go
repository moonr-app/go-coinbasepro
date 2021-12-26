package coinbasepro

import (
	"context"
	"fmt"
	"net/http"
)

type Currency struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	MinSize string `json:"min_size"`
}

func (c *client) GetCurrencies(ctx context.Context) ([]Currency, error) {
	var currencies []Currency

	url := fmt.Sprintf("/currencies")
	_, err := c.Request(ctx, http.MethodGet, url, nil, &currencies)
	return currencies, err
}
