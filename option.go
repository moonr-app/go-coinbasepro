package coinbasepro

import (
	"errors"
	"net/http"
	"time"
)

func WithSandboxEnvironment() ClientOption {
	return func(c *client) error {
		c.baseURL = baseURLSandbox
		c.websocketURL = websocketURLSandbox
		return nil
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *client) error {
		if httpClient == nil {
			return errors.New("httpClient cannot be nil")
		}
		c.httpClient = httpClient

		return nil
	}
}

func WithRetryCount(retryCount int) ClientOption {
	return func(c *client) error {
		if retryCount < 0 {
			return errors.New("retryCount cannot be less than 0")
		}
		c.retryCount = retryCount

		return nil
	}
}

func WithRetryInterval(retryInterval time.Duration) ClientOption {
	return func(c *client) error {
		if retryInterval < 0 {
			return errors.New("retryInterval cannot be less than 0")
		}
		c.retryInterval = retryInterval

		return nil
	}
}

// WithTimeOffsetSeconds can be used to generate timestamps with offset to current time.
// Coinbase Pro sandbox time has been off in the past.
func WithTimeOffsetSeconds(offset int) ClientOption {
	return func(c *client) error {
		c.timeOffsetSeconds = offset

		return nil
	}
}
