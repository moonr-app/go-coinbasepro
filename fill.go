package coinbasepro

import (
	"fmt"
	"net/http"
)

type Fill struct {
	TradeID   int    `json:"trade_id,int"`
	ProductID string `json:"product_id"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	FillID    string `json:"order_id"`
	CreatedAt Time   `json:"created_at,string"`
	Fee       string `json:"fee"`
	Settled   bool   `json:"settled"`
	Side      string `json:"side"`
	Liquidity string `json:"liquidity"`
}

type ListFillsParams struct {
	OrderID    string
	ProductID  string
	Pagination PaginationParams
}

func (c *client) ListFills(p ListFillsParams) *Cursor {
	paginationParams := p.Pagination
	if p.OrderID != "" {
		paginationParams.AddExtraParam("order_id", p.OrderID)
	}
	if p.ProductID != "" {
		paginationParams.AddExtraParam("product_id", p.ProductID)
	}

	return c.newCursor(http.MethodGet, fmt.Sprintf("/fills"), paginationParams)
}
