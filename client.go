package coinbasepro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	baseURLProduction = "https://api.pro.coinbase.com"
	baseURLSandbox    = "https://api-public.sandbox.pro.coinbase.com"
)

type (
	client struct {
		baseURL           string
		secret            string
		key               string
		passphrase        string
		httpClient        *http.Client
		retryCount        int
		retryInterval     time.Duration
		timeOffsetSeconds int
	}
	option func(*client) error
)

func NewClient(key, passphrase, secret string, opts ...option) (*client, error) {
	c := &client{
		baseURL:           baseURLProduction,
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

func (c *client) Request(method string, url string,
	params, result interface{}) (res *http.Response, err error) {
	for i := 0; i < c.retryCount+1; i++ {
		retryDuration := time.Duration((math.Pow(2, float64(i))-1)/2) * c.retryInterval
		time.Sleep(retryDuration)

		res, err = c.request(method, url, params, result)
		if res != nil && res.StatusCode == 429 {
			continue
		}

		break
	}

	return res, err
}

func (c *client) request(method string, url string, params, result interface{}) (res *http.Response, err error) {
	var data []byte
	body := bytes.NewReader(make([]byte, 0))

	if params != nil {
		data, err = json.Marshal(params)
		if err != nil {
			return res, err
		}

		body = bytes.NewReader(data)
	}

	fullURL := fmt.Sprintf("%s%s", c.baseURL, url)
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return res, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix()+int64(c.timeOffsetSeconds), 10)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Go Coinbase Pro httpClient 1.0")

	h, err := c.Headers(method, url, timestamp, string(data))
	if err != nil {
		return res, err
	}

	for k, v := range h {
		req.Header.Add(k, v)
	}

	res, err = c.httpClient.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		coinbaseError := Error{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&coinbaseError); err != nil {
			return res, err
		}

		return res, error(coinbaseError)
	}

	if result != nil {
		decoder := json.NewDecoder(res.Body)
		if err = decoder.Decode(result); err != nil {
			return res, err
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
		return nil, err
	}
	h["CB-ACCESS-SIGN"] = sig
	return h, nil
}
