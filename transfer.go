package coinbasepro

import (
	"context"
	"fmt"
	"net/http"
)

type Transfer struct {
	Type              string `json:"type"`
	Amount            string `json:"amount"`
	CoinbaseAccountID string `json:"coinbase_account_id,string"`
}

func (c *client) CreateTransfer(ctx context.Context, newTransfer Transfer) (Transfer, error) {
	var savedTransfer Transfer

	url := fmt.Sprintf("/transfers")
	_, err := c.Request(ctx, http.MethodPost, url, newTransfer, &savedTransfer)
	return savedTransfer, err
}
