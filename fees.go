package coinbasepro

import (
	"context"
	"fmt"
	"net/http"
)

type Fees struct {
	MakerFeeRate string `json:"maker_fee_rate"`
	TakerFeeRate string `json:"taker_fee_rate"`
	USDVolume    string `json:"usd_volume"`
}

func (c *client) GetFees(ctx context.Context) (Fees, error) {
	var fees Fees

	url := fmt.Sprintf("/fees")
	_, err := c.Request(ctx, http.MethodGet, url, nil, &fees)
	return fees, err
}
