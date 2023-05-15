# SMS Partner

## Summary

SMSPartner API client in Go

## Requirements:

To use client v1, you'll need to set the environment variable

- `SMSPARTNER_API_KEY`

## Installation

```go
$ go get -u -v github.com/hoflish/smspartner-go
```

## Documentation

[API Documentation (fr)](https://my.smspartner.fr/documentation-fr/api/v1)

## Usage

Sample usage: You can see file
[example_test.go](./example_test.go)

### With default client

- Check credits

```go
import (
	"fmt"
	"log"
	"net/http"

	"github.com/hoflish/smspartner-go/v1"
)

  client, err := smspartner.NewClient(&http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	credits, err := client.CheckCredits()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Credits: %#v\n", credits)
```

### With custom client

```go
import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hoflish/smspartner-go/v1"
)


  var tr = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		// ...
	}
	var client = &http.Client{
		Transport: tr,
	}
	spClient, err := smspartner.NewClient(client)
	if err != nil {
		log.Fatal(err)
	}

	credits, err := spClient.CheckCredits()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Credits: %#v\n", credits)
```

## Test

Run all tests:

    go test ./...

Run tests for one package:

    go test ./v1

Run a single test:

    go test ./v1 -run TestCredits

