package coinbasepro

import (
	"context"
	"fmt"
	"net/http"
)

type Account struct {
	ID        string `json:"id"`
	Balance   string `json:"balance"`
	Hold      string `json:"hold"`
	Available string `json:"available"`
	Currency  string `json:"currency"`
}

// Ledger

type LedgerEntry struct {
	ID        string        `json:"id,number"`
	CreatedAt Time          `json:"created_at,string"`
	Amount    string        `json:"amount"`
	Balance   string        `json:"balance"`
	Type      string        `json:"type"`
	Details   LedgerDetails `json:"details"`
}

type LedgerDetails struct {
	OrderID   string `json:"order_id"`
	TradeID   string `json:"trade_id"`
	ProductID string `json:"product_id"`
}

type GetAccountLedgerParams struct {
	Pagination PaginationParams
}

// Holds

type Hold struct {
	AccountID string `json:"account_id"`
	CreatedAt Time   `json:"created_at,string"`
	UpdatedAt Time   `json:"updated_at,string"`
	Amount    string `json:"amount"`
	Type      string `json:"type"`
	Ref       string `json:"ref"`
}

type ListHoldsParams struct {
	Pagination PaginationParams
}

func (c *client) GetAccounts(ctx context.Context) ([]Account, error) {
	var accounts []Account
	_, err := c.Request(ctx, http.MethodGet, "/accounts", nil, &accounts)

	return accounts, err
}

func (c *client) GetAccount(ctx context.Context, id string) (Account, error) {
	account := Account{}

	url := fmt.Sprintf("/accounts/%s", id)
	_, err := c.Request(ctx, http.MethodGet, url, nil, &account)
	return account, err
}

func (c *client) ListAccountLedger(id string, p ...GetAccountLedgerParams) *Cursor {
	paginationParams := PaginationParams{}
	if len(p) > 0 {
		paginationParams = p[0].Pagination
	}

	return c.newCursor(http.MethodGet, fmt.Sprintf("/accounts/%s/ledger", id), paginationParams)
}

func (c *client) ListHolds(id string, p ...ListHoldsParams) *Cursor {
	paginationParams := PaginationParams{}
	if len(p) > 0 {
		paginationParams = p[0].Pagination
	}

	return c.newCursor(http.MethodGet, fmt.Sprintf("/accounts/%s/holds", id), paginationParams)
}
