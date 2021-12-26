package coinbasepro

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	ws "github.com/gorilla/websocket"
)

func NewTestClient(t *testing.T) *client {
	c, err := NewClient(
		os.Getenv("COINBASE_PRO_KEY"),
		os.Getenv("COINBASE_PRO_PASSPHRASE"),
		os.Getenv("COINBASE_PRO_SECRET"),
		WithSandboxEnvironment(),
		WithRetryCount(2),
	)
	if err != nil {
		t.Fatal(err)
	}

	return c
}

func NewTestWebsocketClient() (*ws.Conn, error) {
	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed-public.sandbox.pro.coinbase.com", nil)

	return wsConn, err
}

func StructHasZeroValues(i interface{}) bool {
	iv := reflect.ValueOf(i)

	for i := 0; i < iv.NumField(); i++ {
		field := iv.Field(i)
		if reflect.Zero(field.Type()).Interface() == field.Interface() {
			return true
		}
	}

	return false
}

func CompareProperties(a, b interface{}, properties []string) (bool, error) {
	aValueOf := reflect.ValueOf(a)
	bValueOf := reflect.ValueOf(b)

	for _, property := range properties {
		aValue := reflect.Indirect(aValueOf).FieldByName(property).Interface()
		bValue := reflect.Indirect(bValueOf).FieldByName(property).Interface()

		if aValue != bValue {
			return false, fmt.Errorf(fmt.Sprintf("%s not equal: %s - %s", property, aValue, bValue))
		}
	}

	return true, nil
}

func Ensure(a interface{}) error {
	field := reflect.Indirect(reflect.ValueOf(a))

	switch field.Kind() {
	case reflect.Slice:
		if reflect.ValueOf(field.Interface()).Len() == 0 {
			return fmt.Errorf(fmt.Sprintf("Slice is zero"))
		}
	default:
		if reflect.Zero(field.Type()).Interface() == field.Interface() {
			return fmt.Errorf(fmt.Sprintf("Property is zero"))
		}
	}

	return nil
}

func EnsureProperties(a interface{}, properties []string) error {
	valueOf := reflect.ValueOf(a)

	for _, property := range properties {
		field := reflect.Indirect(valueOf).FieldByName(property)

		if err := Ensure(field.Interface()); err != nil {
			return fmt.Errorf(fmt.Sprintf("%s: %s", err.Error(), property))
		}
	}

	return nil
}
