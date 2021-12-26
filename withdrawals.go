package coinbasepro

import (
	"context"
	"fmt"
	"net/http"
)

type WithdrawalCrypto struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	CryptoAddress string `json:"crypto_address"`
}

type WithdrawalCoinbase struct {
	Currency          string `json:"currency"`
	Amount            string `json:"amount"`
	CoinbaseAccountID string `json:"coinbase_account_id"`
}

type WithdrawalPaymentMethod struct {
	ProfileID string `json:"profile_id"`
	Currency  string `json:"currency"`
	Amount    string `json:"amount"`
	// PaymentMethodID can be determined by calling GetPaymentMethods
	PaymentMethodID string `json:"payment_method_id"`
	// Response fields
	ID       string `json:"id,omitempty"`
	PayoutAt Time   `json:"payout_at,string,omitempty"`
	Fee      string `json:"fee,omitempty"`
	Subtotal string `json:"subtotal,omitempty"`
}

func (c *client) CreateWithdrawalPaymentMethod(ctx context.Context, newWithdrawal WithdrawalPaymentMethod) (WithdrawalPaymentMethod, error) {
	var savedWithdrawal WithdrawalPaymentMethod

	url := fmt.Sprintf("/withdrawals/payment-method")
	_, err := c.Request(ctx, http.MethodPost, url, newWithdrawal, &savedWithdrawal)
	return savedWithdrawal, err
}

func (c *client) CreateWithdrawalCrypto(ctx context.Context, newWithdrawalCrypto WithdrawalCrypto) (WithdrawalCrypto, error) {
	var savedWithdrawal WithdrawalCrypto
	url := fmt.Sprintf("/withdrawals/crypto")
	_, err := c.Request(ctx, http.MethodPost, url, newWithdrawalCrypto, &savedWithdrawal)
	return savedWithdrawal, err
}

func (c *client) CreateWithdrawalCoinbase(ctx context.Context, newWithdrawalCoinbase WithdrawalCoinbase) (WithdrawalCoinbase, error) {
	var savedWithdrawal WithdrawalCoinbase
	url := fmt.Sprintf("/withdrawals/coinbase-account")
	_, err := c.Request(ctx, http.MethodPost, url, newWithdrawalCoinbase, &savedWithdrawal)
	return savedWithdrawal, err
}
