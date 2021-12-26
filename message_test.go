package coinbasepro_test

import (
	"context"
	"testing"

	"github.com/preichenberger/go-coinbasepro/v2"
)

func TestMessageHeartbeat(t *testing.T) {
	client := coinbasepro.NewTestClient(t)

	subscribe := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			{
				Name: "heartbeat",
				ProductIds: []string{
					"BTC-USD",
				},
			},
		},
	}

	var (
		message coinbasepro.Message
		ctx     = context.Background()
	)

	err := client.Subscribe(ctx, subscribe, func(msg coinbasepro.Message) error {
		message = msg
		if msg.Type != "subscriptions" {
			return coinbasepro.ErrCloseWebsocket
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if message.Type != "heartbeat" {
		t.Fatal("Invalid message type")
	}

	// LastTradeId is broken on sandbox
	// props := []string{"Type", "Sequence", "LastTradeId", "ProductId", "Time"}"
	props := []string{"Type", "Sequence", "ProductID", "Time"}
	if err := coinbasepro.EnsureProperties(message, props); err != nil {
		t.Fatal(err)
	}
}

func TestMessageTicker(t *testing.T) {
	client := coinbasepro.NewTestClient(t)

	subscribe := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			{
				Name: "ticker",
				ProductIds: []string{
					"BTC-USD",
				},
			},
		},
	}

	var (
		message coinbasepro.Message
		ctx     = context.Background()
	)

	err := client.Subscribe(ctx, subscribe, func(msg coinbasepro.Message) error {
		message = msg
		if msg.Type != "subscriptions" {
			return coinbasepro.ErrCloseWebsocket
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if message.Type != "ticker" {
		t.Fatal("Invalid message type")
	}

	props := []string{"Type", "Sequence", "ProductID", "BestBid", "BestAsk", "Price"}
	if err := coinbasepro.EnsureProperties(message, props); err != nil {
		t.Fatal(err)
	}
}

func TestMessageLevel2(t *testing.T) {
	client := coinbasepro.NewTestClient(t)

	subscribe := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			{
				Name: "level2",
				ProductIds: []string{
					"BTC-USD",
				},
			},
		},
	}

	var (
		message *coinbasepro.Message
		ctx     = context.Background()
		i       = 0
		l2      = false
	)

	err := client.Subscribe(ctx, subscribe, func(msg coinbasepro.Message) error {
		if message == nil && msg.Type != "subscriptions" {
			message = &msg
		}

		if msg.Type == "l2update" {
			l2 = true
			props := []string{"ProductID", "Changes"}
			if err := coinbasepro.EnsureProperties(msg, props); err != nil {
				t.Fatal(err)
			}

			return coinbasepro.ErrCloseWebsocket
		}

		if i == 10 {
			return coinbasepro.ErrCloseWebsocket
		}
		i++

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if message.Type != "snapshot" {
		t.Fatal("Invalid message type")
	}

	props := []string{"ProductID", "Bids", "Asks"}
	if err := coinbasepro.EnsureProperties(message, props); err != nil {
		t.Fatal(err)
	}

	if !l2 {
		t.Fatal("Did not find l2update")
	}
}

func TestMessageStatus(t *testing.T) {
	client := coinbasepro.NewTestClient(t)

	subscribe := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			{
				Name: "status",
			},
		},
	}

	var (
		message coinbasepro.Message
		ctx     = context.Background()
	)

	err := client.Subscribe(ctx, subscribe, func(msg coinbasepro.Message) error {
		message = msg
		if msg.Type != "subscriptions" {
			return coinbasepro.ErrCloseWebsocket
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if message.Type != "status" {
		t.Fatal("invalid message type")
	}

	props := []string{"Products", "Currencies"}
	if err := coinbasepro.EnsureProperties(message, props); err != nil {
		t.Fatal(err)
	}
}
