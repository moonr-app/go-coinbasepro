package coinbasepro

import (
	"context"
	"fmt"
	"net/http"
)

type Deposit struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
	// PaymentMethodID can be determined by calling GetPaymentMethods
	PaymentMethodID string `json:"payment_method_id"`
	// Response fields
	ID       string `json:"id,omitempty"`
	PayoutAt Time   `json:"payout_at,string,omitempty"`
}

func (c *client) CreateDeposit(ctx context.Context, newDeposit Deposit) (Deposit, error) {
	var savedDeposit Deposit

	url := fmt.Sprintf("/deposits/payment-method")
	_, err := c.Request(ctx, http.MethodPost, url, newDeposit, &savedDeposit)
	return savedDeposit, err
}

type PaymentMethod struct {
	Currency string `json:"currency"`
	Type     string `json:"type"`
	ID       string `json:"id"`
}

func (c *client) GetPaymentMethods(ctx context.Context) ([]PaymentMethod, error) {
	var paymentMethods []PaymentMethod

	url := fmt.Sprintf("/payment-methods")
	_, err := c.Request(ctx, http.MethodGet, url, nil, &paymentMethods)

	return paymentMethods, err
}
