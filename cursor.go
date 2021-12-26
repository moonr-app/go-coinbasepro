package coinbasepro

import (
	"context"
	"fmt"
)

type Cursor struct {
	client     *client
	pagination PaginationParams
	method     string
	params     interface{}
	url        string
	HasMore    bool
}

func (c *client) newCursor(method, url string, paginationParams PaginationParams) *Cursor {
	return &Cursor{
		client:     c,
		method:     method,
		url:        url,
		pagination: paginationParams,
		HasMore:    true,
	}
}

func (c *Cursor) page(ctx context.Context, i interface{}, direction string) error {
	url := c.url
	if c.pagination.Encode(direction) != "" {
		url = fmt.Sprintf("%s?%s", c.url, c.pagination.Encode(direction))
	}

	res, err := c.client.Request(ctx, c.method, url, c.params, i)
	if err != nil {
		c.HasMore = false
		return err
	}

	c.pagination.Before = res.Header.Get("CB-BEFORE")
	c.pagination.After = res.Header.Get("CB-AFTER")

	if c.pagination.Done(direction) {
		c.HasMore = false
	}

	return nil
}

func (c *Cursor) NextPage(ctx context.Context, i interface{}) error {
	return c.page(ctx, i, "next")
}

func (c *Cursor) PrevPage(ctx context.Context, i interface{}) error {
	return c.page(ctx, i, "prev")
}
