package coinbasepro_test

import (
	"context"
	"testing"

	"github.com/moonr-app/go-coinbasepro"
)

func TestCreateLimitOrders(t *testing.T) {
	client := coinbasepro.NewTestClient(t)

	order := coinbasepro.Order{
		Price:     "1.00000000",
		Size:      "2.00000000",
		Side:      "buy",
		ProductID: "BTC-GBP",
	}
	ctx := context.Background()

	savedOrder, err := client.CreateOrder(ctx, order)
	if err != nil {
		t.Fatal(err)
	}

	if savedOrder.ID == "" {
		t.Fatal("No create id found")
	}

	props := []string{"Price", "Size", "Side", "ProductID"}
	_, err = coinbasepro.CompareProperties(order, savedOrder, props)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.CancelOrder(ctx, savedOrder.ID); err != nil {
		t.Fatal(err)
	}
}

func TestCreateMarketOrders(t *testing.T) {
	client := coinbasepro.NewTestClient(t)

	order := coinbasepro.Order{
		Funds:     "10.00",
		Size:      "1.50000000",
		Side:      "buy",
		Type:      "market",
		ProductID: "BTC-GBP",
	}
	ctx := context.Background()

	savedOrder, err := client.CreateOrder(ctx, order)
	if err != nil {
		t.Fatal(err)
	}

	if savedOrder.ID == "" {
		t.Fatal("No create id found")
	}

	props := []string{"Price", "Size", "Side", "ProductID"}
	_, err = coinbasepro.CompareProperties(order, savedOrder, props)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCancelOrder(t *testing.T) {
	client := coinbasepro.NewTestClient(t)

	order := coinbasepro.Order{
		Price:     "1.00",
		Size:      "1.30",
		Side:      "buy",
		ProductID: "BTC-GBP",
	}
	ctx := context.Background()

	savedOrder, err := client.CreateOrder(ctx, order)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.CancelOrder(ctx, savedOrder.ID); err != nil {
		t.Fatal(err)
	}
}

func TestGetOrder(t *testing.T) {
	client := coinbasepro.NewTestClient(t)

	order := coinbasepro.Order{
		Price:     "1.00",
		Size:      "1.00",
		Side:      "buy",
		ProductID: "BTC-GBP",
	}
	ctx := context.Background()

	savedOrder, err := client.CreateOrder(ctx, order)
	if err != nil {
		t.Fatal(err)
	}

	getOrder, err := client.GetOrder(ctx, savedOrder.ID)
	if err != nil {
		t.Fatal(err)
	}

	if getOrder.ID != savedOrder.ID {
		t.Fatal("Order ids do not match")
	}

	if err := client.CancelOrder(ctx, savedOrder.ID); err != nil {
		t.Fatal(err)
	}
}

func TestListOrders(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	ctx := context.Background()

	cursor := client.ListOrders(coinbasepro.ListOrdersParams{})
	var orders []coinbasepro.Order

	for cursor.HasMore {
		if err := cursor.NextPage(ctx, &orders); err != nil {
			t.Fatal(err)
		}

		for _, o := range orders {
			if coinbasepro.StructHasZeroValues(o) {
				t.Fatal("Zero value")
			}
		}
	}

	cursor = client.ListOrders(coinbasepro.ListOrdersParams{Status: "open", ProductID: "BTC-EUR"})
	for cursor.HasMore {
		if err := cursor.NextPage(ctx, &orders); err != nil {
			t.Fatal(err)
		}

		for _, o := range orders {
			if coinbasepro.StructHasZeroValues(o) {
				t.Fatal("Zero value")
			}
		}
	}
}

func TestCancelAllOrders(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	ctx := context.Background()

	for _, pair := range []string{"BTC-GBP"} {
		for i := 0; i < 2; i++ {
			order := coinbasepro.Order{Price: "100.00", Size: "1.00", Side: "buy", ProductID: pair}

			if _, err := client.CreateOrder(ctx, order); err != nil {
				t.Fatal(err)
			}
		}

		orderIDs, err := client.CancelAllOrders(ctx, coinbasepro.CancelAllOrdersParams{ProductID: pair})
		if err != nil {
			t.Fatal(err)
		}

		if len(orderIDs) != 2 {
			t.Fatal("Did not cancel all orders")
		}
	}
}
