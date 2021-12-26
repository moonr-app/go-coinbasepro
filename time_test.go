package coinbasepro_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/moonr-app/go-coinbasepro"
)

func TestGetTime(t *testing.T) {
	client := coinbasepro.NewTestClient(t)
	serverTime, err := client.GetTime(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if coinbasepro.StructHasZeroValues(serverTime) {
		t.Fatal("Zero value")
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	c := coinbasepro.Time{}
	now := time.Now()

	jsonData, err := json.Marshal(now.Format("2006-01-02 15:04:05+00"))
	if err != nil {
		t.Fatal(err)
	}

	if err = c.UnmarshalJSON(jsonData); err != nil {
		t.Fatal(err)
	}

	if now.Equal(c.Time()) {
		t.Fatal("Unmarshaled time does not equal original time")
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	c := coinbasepro.Time{}
	tt := time.Date(9999, 4, 12, 23, 20, 50, 0, time.UTC)
	expected := "\"9999-04-12T23:20:50Z\""

	jsonData, err := json.Marshal(tt.Format("2006-01-02 15:04:05+00"))
	if err != nil {
		t.Fatal(err)
	}

	if err = c.UnmarshalJSON(jsonData); err != nil {
		t.Fatal(err)
	}

	jsonData, err = json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}

	if string(jsonData) != expected {
		t.Fatal("Marshaled time (" + string(jsonData) + ") does not equal original time (" + expected + ")")
	}
}
