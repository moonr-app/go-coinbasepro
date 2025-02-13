package coinbasepro

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	ws "github.com/gorilla/websocket"
)

type Message struct {
	Type          string           `json:"type"`
	ProductID     string           `json:"product_id"`
	ProductIds    []string         `json:"product_ids"`
	Products      []Product        `json:"products"`
	Currencies    []Currency       `json:"currencies"`
	TradeID       int              `json:"trade_id,number"`
	OrderID       string           `json:"order_id"`
	ClientOID     string           `json:"client_oid"`
	Sequence      int64            `json:"sequence,number"`
	MakerOrderID  string           `json:"maker_order_id"`
	TakerOrderID  string           `json:"taker_order_id"`
	Time          Time             `json:"time,string"`
	RemainingSize string           `json:"remaining_size"`
	NewSize       string           `json:"new_size"`
	OldSize       string           `json:"old_size"`
	Size          string           `json:"size"`
	Price         string           `json:"price"`
	Side          string           `json:"side"`
	Reason        string           `json:"reason"`
	OrderType     string           `json:"order_type"`
	Funds         string           `json:"funds"`
	NewFunds      string           `json:"new_funds"`
	OldFunds      string           `json:"old_funds"`
	Message       string           `json:"message"`
	Bids          []SnapshotEntry  `json:"bids,omitempty"`
	Asks          []SnapshotEntry  `json:"asks,omitempty"`
	Changes       []SnapshotChange `json:"changes,omitempty"`
	LastSize      string           `json:"last_size"`
	BestBid       string           `json:"best_bid"`
	BestAsk       string           `json:"best_ask"`
	Channels      []MessageChannel `json:"channels"`
	UserID        string           `json:"user_id"`
	ProfileID     string           `json:"profile_id"`
	LastTradeID   int              `json:"last_trade_id"`
}

type MessageChannel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

type SnapshotChange struct {
	Side  string
	Price string
	Size  string
}

type SnapshotEntry struct {
	Price string
	Size  string
}

type SignedMessage struct {
	Message
	Key        string `json:"key"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Signature  string `json:"signature"`
}

func (c *client) Subscribe(ctx context.Context, message Message, handler func(Message) error) error {
	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.DialContext(ctx, c.websocketURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}
	defer wsConn.Close()

	if err := wsConn.WriteJSON(message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	for {
		var receivedMessage Message
		if err := wsConn.ReadJSON(&receivedMessage); err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		err := handler(receivedMessage)
		if err == nil {
			continue
		}

		if !errors.Is(err, ErrCloseWebsocket) {
			return fmt.Errorf("failed to handle message: %w", err)
		}

		return nil
	}
}

func (e *SnapshotEntry) UnmarshalJSON(data []byte) error {
	var entry []string

	if err := json.Unmarshal(data, &entry); err != nil {
		return err
	}

	e.Price = entry[0]
	e.Size = entry[1]

	return nil
}

func (e *SnapshotChange) UnmarshalJSON(data []byte) error {
	var entry []string

	if err := json.Unmarshal(data, &entry); err != nil {
		return err
	}

	e.Side = entry[0]
	e.Price = entry[1]
	e.Size = entry[2]

	return nil
}
