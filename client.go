package coinbasepro

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	baseURLProduction      = "https://api.exchange.coinbase.com"
	baseURLSandbox         = "https://api-public.sandbox.exchange.coinbase.com"
	websocketURLProduction = "wss://ws-feed.exchange.coinbase.com"
	websocketURLSandbox    = "wss://ws-feed-public.sandbox.exchange.coinbase.com"
)

type (
	client struct {
		baseURL           string
		websocketURL      string
		secret            string
		key               string
		passphrase        string
		httpClient        *http.Client
		retryCount        int
		retryInterval     time.Duration
		timeOffsetSeconds int
	}
	ClientOption func(*client) error
)

// NewAnonymousClient creates a new instance of client without any credentials which can be used for public endpoints.
func NewAnonymousClient(opts ...ClientOption) (*client, error) {
	c := &client{
		baseURL:           baseURLProduction,
		websocketURL:      websocketURLProduction,
		httpClient:        &http.Client{Timeout: 15 * time.Second},
		retryCount:        0,
		retryInterval:     100 * time.Millisecond,
		timeOffsetSeconds: 0,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// NewClient creates a new instance of client with credentials which can be used for both public & private endpoints.
func NewClient(key, passphrase, secret string, opts ...ClientOption) (*client, error) {
	switch {
	case key == "":
		return nil, errors.New("key cannot be empty")
	case passphrase == "":
		return nil, errors.New("passphrase cannot be empty")
	case secret == "":
		return nil, errors.New("secret cannot be empty")
	}
	c := &client{
		baseURL:           baseURLProduction,
		websocketURL:      websocketURLProduction,
		key:               key,
		passphrase:        passphrase,
		secret:            secret,
		httpClient:        &http.Client{Timeout: 15 * time.Second},
		retryCount:        0,
		retryInterval:     100 * time.Millisecond,
		timeOffsetSeconds: 0,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *client) Request(ctx context.Context, method, url string, params, result interface{}) (res *http.Response, err error) {
	for i := 0; i < c.retryCount+1; i++ {
		retryDuration := time.Duration((math.Pow(2, float64(i))-1)/2) * c.retryInterval
		time.Sleep(retryDuration)

		res, err = c.request(ctx, method, url, params, result)
		if res != nil && res.StatusCode == 429 {
			continue
		}

		break
	}

	return res, err
}

func (c *client) request(ctx context.Context, method, url string, params, result interface{}) (res *http.Response, err error) {
	var data []byte
	body := bytes.NewReader(make([]byte, 0))

	if params != nil {
		data, err = json.Marshal(params)
		if err != nil {
			return res, fmt.Errorf("failed to marshal request: %w", err)
		}

		body = bytes.NewReader(data)
	}

	fullURL := fmt.Sprintf("%s%s", c.baseURL, url)
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return res, fmt.Errorf("failed to create new request: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix()+int64(c.timeOffsetSeconds), 10)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Go Coinbase Pro httpClient 1.0")

	h, err := c.Headers(method, url, timestamp, string(data))
	if err != nil {
		return res, fmt.Errorf("failed to generate request headers: %w", err)
	}

	for k, v := range h {
		req.Header.Add(k, v)
	}

	res, err = c.httpClient.Do(req)
	if err != nil {
		return res, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		coinbaseErr := Error{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&coinbaseErr); err != nil {
			return res, fmt.Errorf("failed to decode error response: %w", err)
		}

		return res, coinbaseErr
	}

	if result != nil {
		decoder := json.NewDecoder(res.Body)
		if err = decoder.Decode(result); err != nil {
			return res, fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return res, nil
}

// Headers generates a map that can be used as headers to authenticate a request
func (c *client) Headers(method, url, timestamp, data string) (map[string]string, error) {
	h := make(map[string]string)
	h["CB-ACCESS-KEY"] = c.key
	h["CB-ACCESS-PASSPHRASE"] = c.passphrase
	h["CB-ACCESS-TIMESTAMP"] = timestamp

	message := fmt.Sprintf(
		"%s%s%s%s",
		timestamp,
		method,
		url,
		data,
	)

	sig, err := generateSig(message, c.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signature: %w", err)
	}
	h["CB-ACCESS-SIGN"] = sig
	return h, nil
}
