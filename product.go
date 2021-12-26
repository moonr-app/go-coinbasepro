package coinbasepro

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Product struct {
	ID              string `json:"id"`
	BaseCurrency    string `json:"base_currency"`
	QuoteCurrency   string `json:"quote_currency"`
	BaseMinSize     string `json:"base_min_size"`
	BaseMaxSize     string `json:"base_max_size"`
	QuoteIncrement  string `json:"quote_increment"`
	BaseIncrement   string `json:"base_increment"`
	DisplayName     string `json:"display_name"`
	MinMarketFunds  string `json:"min_market_funds"`
	MaxMarketFunds  string `json:"max_market_funds"`
	MarginEnabled   bool   `json:"margin_enabled"`
	PostOnly        bool   `json:"post_only"`
	LimitOnly       bool   `json:"limit_only"`
	CancelOnly      bool   `json:"cancel_only"`
	TradingDisabled bool   `json:"trading_disabled"`
	Status          string `json:"status"`
	StatusMessage   string `json:"status_message"`
}

type Ticker struct {
	TradeID int          `json:"trade_id,number"`
	Price   string       `json:"price"`
	Size    string       `json:"size"`
	Time    Time         `json:"time,string"`
	Bid     string       `json:"bid"`
	Ask     string       `json:"ask"`
	Volume  StringNumber `json:"volume"`
}

type Trade struct {
	TradeID int    `json:"trade_id,number"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Time    Time   `json:"time,string"`
	Side    string `json:"side"`
}

type HistoricRate struct {
	Time   time.Time
	Low    float64
	High   float64
	Open   float64
	Close  float64
	Volume float64
}

type Stats struct {
	Low         string `json:"low"`
	High        string `json:"high"`
	Open        string `json:"open"`
	Volume      string `json:"volume"`
	Last        string `json:"last"`
	Volume30Day string `json:"volume_30day"`
}

type BookEntry struct {
	Price          string
	Size           string
	NumberOfOrders int
	OrderID        string
}

type Book struct {
	Sequence int64       `json:"sequence"`
	Bids     []BookEntry `json:"bids"`
	Asks     []BookEntry `json:"asks"`
}

type ListTradesParams struct {
	Pagination *PaginationParams
}

type GetHistoricRatesParams struct {
	Start       time.Time
	End         time.Time
	Granularity int
}

func (e *BookEntry) UnmarshalJSON(data []byte) error {
	var entry []interface{}

	if err := json.Unmarshal(data, &entry); err != nil {
		return err
	}

	priceString, ok := entry[0].(string)
	if !ok {
		return errors.New("Expected string")
	}

	sizeString, ok := entry[1].(string)
	if !ok {
		return errors.New("Expected string")
	}

	*e = BookEntry{
		Price: priceString,
		Size:  sizeString,
	}

	if numberOfOrdersInt, ok := entry[2].(float64); ok {
		e.NumberOfOrders = int(numberOfOrdersInt)
	} else if orderID, ok := entry[2].(string); ok {
		e.OrderID = orderID
	} else {
		return errors.New("Could not parse 3rd column, tried float64 and string")
	}

	return nil
}

func (e *HistoricRate) UnmarshalJSON(data []byte) error {
	var entry []interface{}

	if err := json.Unmarshal(data, &entry); err != nil {
		return err
	}

	t, ok := entry[0].(float64)
	if !ok {
		return errors.New("Expected float64")
	}

	low, ok := entry[1].(float64)
	if !ok {
		return errors.New("Expected float64")
	}

	high, ok := entry[2].(float64)
	if !ok {
		return errors.New("Expected float64")
	}

	open, ok := entry[3].(float64)
	if !ok {
		return errors.New("Expected float64")
	}

	close, ok := entry[4].(float64)
	if !ok {
		return errors.New("Expected float64")
	}

	volume, ok := entry[5].(float64)
	if !ok {
		return errors.New("Expected float64")
	}

	*e = HistoricRate{
		Time:   time.Unix(int64(t), 0),
		Low:    low,
		High:   high,
		Open:   open,
		Close:  close,
		Volume: volume,
	}

	return nil
}

func (c *client) GetBook(ctx context.Context, product string, level int) (Book, error) {
	var book Book

	requestURL := fmt.Sprintf("/products/%s/book?level=%d", product, level)
	_, err := c.Request(ctx, http.MethodGet, requestURL, nil, &book)
	return book, err
}

func (c *client) GetTicker(ctx context.Context, product string) (Ticker, error) {
	var ticker Ticker

	requestURL := fmt.Sprintf("/products/%s/ticker", product)
	_, err := c.Request(ctx, http.MethodGet, requestURL, nil, &ticker)
	return ticker, err
}

func (c *client) ListTrades(product string, p ListTradesParams) *Cursor {
	paginationParams := PaginationParams{}
	if p.Pagination != nil {
		paginationParams = *p.Pagination
	}

	return c.newCursor(http.MethodGet, fmt.Sprintf("/products/%s/trades", product), paginationParams)
}

func (c *client) GetProducts(ctx context.Context) ([]Product, error) {
	var products []Product

	requestURL := fmt.Sprintf("/products")
	_, err := c.Request(ctx, http.MethodGet, requestURL, nil, &products)
	return products, err
}

func (c *client) GetProduct(ctx context.Context, p string) (Product, error) {
	var product Product
	requestURL := fmt.Sprintf("/products/%s", p)
	_, err := c.Request(ctx, http.MethodGet, requestURL, nil, &product)
	return product, err
}

func (c *client) GetHistoricRates(ctx context.Context, product string, p GetHistoricRatesParams) ([]HistoricRate, error) {
	var historicRates []HistoricRate
	requestURL := fmt.Sprintf("/products/%s/candles", product)
	values := url.Values{}
	layout := "2006-01-02T15:04:05Z"

	if !p.Start.IsZero() {
		values.Add("start", p.Start.UTC().Format(layout))
	}

	if !p.End.IsZero() {
		values.Add("end", p.End.UTC().Format(layout))
	}

	if p.Granularity != 0 {
		values.Add("granularity", strconv.Itoa(p.Granularity))
	}

	if len(values) > 0 {
		requestURL = fmt.Sprintf("%s?%s", requestURL, values.Encode())
	}

	_, err := c.Request(ctx, http.MethodGet, requestURL, nil, &historicRates)
	return historicRates, err
}

func (c *client) GetStats(ctx context.Context, product string) (Stats, error) {
	var stats Stats
	requestURL := fmt.Sprintf("/products/%s/stats", product)
	_, err := c.Request(ctx, http.MethodGet, requestURL, nil, &stats)
	return stats, err
}
