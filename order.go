package coinbasepro

import (
	"context"
	"fmt"
	"net/http"
)

type Order struct {
	Type      string `json:"type"`
	Size      string `json:"size,omitempty"`
	Side      string `json:"side"`
	ProductID string `json:"product_id"`
	ClientOID string `json:"client_oid,omitempty"`
	Stp       string `json:"stp,omitempty"`
	Stop      string `json:"stop,omitempty"`
	StopPrice string `json:"stop_price,omitempty"`
	// Limit Order
	Price       string `json:"price,omitempty"`
	TimeInForce string `json:"time_in_force,omitempty"`
	PostOnly    bool   `json:"post_only,omitempty"`
	CancelAfter string `json:"cancel_after,omitempty"`
	// Market Order
	Funds          string `json:"funds,omitempty"`
	SpecifiedFunds string `json:"specified_funds,omitempty"`
	// Response Fields
	ID            string `json:"id"`
	Status        string `json:"status,omitempty"`
	Settled       bool   `json:"settled,omitempty"`
	DoneReason    string `json:"done_reason,omitempty"`
	DoneAt        Time   `json:"done_at,string,omitempty"`
	CreatedAt     Time   `json:"created_at,string,omitempty"`
	FillFees      string `json:"fill_fees,omitempty"`
	FilledSize    string `json:"filled_size,omitempty"`
	ExecutedValue string `json:"executed_value,omitempty"`
}

type CancelAllOrdersParams struct {
	ProductID string
}

type ListOrdersParams struct {
	Status     string
	ProductID  string
	Pagination PaginationParams
}

func (c *client) CreateOrder(ctx context.Context, newOrder Order) (Order, error) {
	var savedOrder Order

	if len(newOrder.Type) == 0 {
		newOrder.Type = "limit"
	}

	url := fmt.Sprintf("/orders")
	_, err := c.Request(ctx, http.MethodPost, url, newOrder, &savedOrder)
	return savedOrder, err
}

func (c *client) CancelOrder(ctx context.Context, id string) error {
	url := fmt.Sprintf("/orders/%s", id)
	_, err := c.Request(ctx, http.MethodDelete, url, nil, nil)
	return err
}

func (c *client) CancelAllOrders(ctx context.Context, p CancelAllOrdersParams) ([]string, error) {
	var orderIDs []string
	url := "/orders"

	if p.ProductID != "" {
		url = fmt.Sprintf("%s?product_id=%s", url, p.ProductID)
	}

	_, err := c.Request(ctx, http.MethodDelete, url, nil, &orderIDs)
	return orderIDs, err
}

func (c *client) GetOrder(ctx context.Context, id string) (Order, error) {
	var savedOrder Order

	url := fmt.Sprintf("/orders/%s", id)
	_, err := c.Request(ctx, http.MethodGet, url, nil, &savedOrder)
	return savedOrder, err
}

func (c *client) ListOrders(p ListOrdersParams) *Cursor {
	paginationParams := PaginationParams{}
	paginationParams = p.Pagination

	if p.Status != "" {
		paginationParams.AddExtraParam("status", p.Status)
	}
	if p.ProductID != "" {
		paginationParams.AddExtraParam("product_id", p.ProductID)
	}

	return c.newCursor(http.MethodGet, fmt.Sprintf("/orders"), paginationParams)
}
