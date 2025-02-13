Go Coinbase Pro [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/moonr-app/go-coinbasepro) [![Build Status](https://travis-ci.org/moonr-app/go-coinbasepro.svg?branch=master)](https://travis-ci.org/moonr-app/go-coinbasepro)
========

## Summary

Go client for [Coinbase Pro](https://pro.coinbase.com) formerly known as gdax

## Installation
If using Go modules (Go version >= 11.1) simply import as needed.
```sh
go mod init github.com/yourusername/yourprojectname
```

### Older Go versions
```sh
go get github.com/moonr-app/go-coinbasepro
```

### Significant releases
Use [dep](https://github.com/golang/dep) to install previous releases
```sh
dep ensure --add github.com/moonr-app/go-gdax@0.5.7
```

- 0.5.7, last release before rename package to: coinbasepro
- 0.5, as of 0.5 this library uses strings and is not backwards compatible

## Documentation
For full details on functionality, see [GoDoc](http://godoc.org/github.com/moonr-app/go-coinbasepro) documentation.

### Setup
Client will respect environment variables: COINBASE_PRO_BASEURL, COINBASE_PRO_PASSPHRASE, COINBASE_PRO_KEY, COINBASE_PRO_SECRET by default

```go
import (
  coinbasepro "github.com/moonr-app/go-coinbasepro"
)

client, err := coinbasepro.NewClient(
    "coinbase pro key",
    "coinbase pro passphrase",
    "coinbase pro secret",
)
if err != nil {
    // handle error
}

// optional configuration can be provided using function options:
client, err := coinbasepro.NewClient(
    "coinbase pro key",
    "coinbase pro passphrase",
    "coinbase pro secret",
    coinbasepro.WithSandboxEnvironment(),
    coinbasepro.WithHTTPClient(&http.Client{}),
    coinbasepro.WithRetryCount(5),
    coinbasepro.WithRetryInterval(time.Second),
    coinbasepro.WithTimeOffsetSeconds(30),
)
if err != nil {
    // handle error
}
```

### Sandbox
You can switch to the sandbox env by using the following functional option (will default to production if not provided):

```sh
coinbasepro.WithSandboxEnvironment()
```


### HTTP Settings
You can use a custom http client by using the following functional option:
```go
coinbasepro.WithHTTPClient(&http.Client{}),
```

### Decimals
To manage precision correctly, this library sends all price values as strings. It is recommended to use a decimal library
like Spring's [Decimal](https://github.com/shopspring/decimal) if you are doing any manipulation of prices.

Example:
```go
import (
  "github.com/shopspring/decimal"
)

book, err := client.GetBook("BTC-USD", 1)
if err != nil {
    println(err.Error())  
}

lastPrice, err := decimal.NewFromString(book.Bids[0].Price)
if err != nil {
    println(err.Error())  
}

order := coinbasepro.Order{
  Price: lastPrice.Add(decimal.NewFromFloat(1.00)).String(),
  Size: "2.00",
  Side: "buy",
  ProductID: "BTC-USD",
}

savedOrder, err := client.CreateOrder(&order)
if err != nil {
  println(err.Error())
}

println(savedOrder.ID)
```

### Retry
You can set a retry count & interval which uses exponential backoff using functional options:

`(2^(retry_attempt) - 1) / 2 * retryInterval`

```
coinbasepro.WithRetryCount(3),
coinbasepro.WithRetryInterval(500 * time.Millisecond),

// retry count = 3: 500ms, 1500ms, 3500ms
```

### Cursor
This library uses a cursor pattern so you don't have to keep track of pagination.

```go
var orders []coinbasepro.Order
cursor = client.ListOrders()

for cursor.HasMore {
  if err := cursor.NextPage(ctx, &orders); err != nil {
    println(err.Error())
    return
  }

  for _, o := range orders {
    println(o.ID)
  }
}

```

### Websockets
Listen for websocket messages

```go
  client, err := coinbasepro.NewClient(
    "coinbase pro key", 
	"coinbase pro passphrase", 
	"coinbase pro secret",
  )
  if err != nil { 
    // handle error
  }

  subscribe := coinbasepro.Message{
    Type:      "subscribe",
    Channels: []coinbasepro.MessageChannel{
      coinbasepro.MessageChannel{
        Name: "heartbeat",
        ProductIds: []string{
          "BTC-USD",
        },
      },
      coinbasepro.MessageChannel{
        Name: "level2",
        ProductIds: []string{
          "BTC-USD",
        },
      },
    },
  }
  
  err := client.Subscribe(ctx, subscribe, func(msg coinbasepro.Message) error {
    println(message.Type)
    return nil
  })
```

### Time
Results return coinbase time type which handles different types of time parsing that coinbasepro returns. This wraps the native go time type

```go
  import(
    "time"
    coinbasepro "github.com/moonr-app/go-coinbasepro"
  )

  coinbaseTime := coinbasepro.Time{}
  println(time.Time(coinbaseTime).Day())
```

### Examples
This library supports all public and private endpoints

Get Accounts:
```go
  accounts, err := client.GetAccounts(ctx)
  if err != nil {
    println(err.Error())
  }

  for _, a := range accounts {
    println(a.Balance)
  }
```

List Account Ledger:
```go
  var ledgers []coinbasepro.LedgerEntry

  accounts, err := client.GetAccounts(ctx)
  if err != nil {
    println(err.Error())
  }

  for _, a := range accounts {
    cursor := client.ListAccountLedger(a.ID)
    for cursor.HasMore {
      if err := cursor.NextPage(ctx, &ledgers); err != nil {
        println(err.Error())
      }

      for _, e := range ledgers {
        println(e.Amount)
      }
    }
  }
```

Create an Order:
```go
  order := coinbasepro.Order{
    Price: "1.00",
    Size: "1.00",
    Side: "buy",
    ProductID: "BTC-USD",
  }

  savedOrder, err := client.CreateOrder(ctx, order)
  if err != nil {
    println(err.Error())
  }

  println(savedOrder.ID)
```

Transfer funds:
```go
  transfer := coinbasepro.Transfer {
    Type: "deposit",
    Amount: "1.00",
  }

  savedTransfer, err := client.CreateTransfer(ctx, transfer)
  if err != nil {
    println(err.Error())
  }
```

Get Trade history:
```go
  var trades []coinbasepro.Trade
  cursor := client.ListTrades("BTC-USD")

  for cursor.HasMore {
    if err := cursor.NextPage(ctx, &trades); err != nil {
      for _, t := range trades {
        println(trade.CoinbaseID)
      }
    }
  }
```

### Testing
To test with Coinbase's public sandbox set the following environment variables:
```sh
export COINBASE_PRO_KEY="sandbox key"
export COINBASE_PRO_PASSPHRASE="sandbox passphrase"
export COINBASE_PRO_SECRET="sandbox secret"
```

Then run `go test`
```sh
go test
```

Note that your sandbox account will need at least 2,000 USD and 2 BTC to run the tests.