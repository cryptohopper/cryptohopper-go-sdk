# cryptohopper-go-sdk

Official Go SDK for the [Cryptohopper](https://www.cryptohopper.com) API.

> **Status: v0.1.0-alpha.1** — early access. Covers six core domains: `User`, `Hoppers`, `Exchange`, `Strategy`, `Backtest`, `Market`.

## Install

```bash
go get github.com/cryptohopper/cryptohopper-go-sdk@v0.1.0-alpha.1
```

Requires Go 1.22+.

## Quickstart

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/cryptohopper/cryptohopper-go-sdk"
)

func main() {
    client, err := cryptohopper.NewClient(os.Getenv("CRYPTOHOPPER_TOKEN"))
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    me, err := client.User.Get(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(me["email"])

    ticker, err := client.Exchange.Ticker(ctx, "binance", "BTC/USDT")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(ticker["last"])
}
```

## Authentication

Cryptohopper uses OAuth2 bearer tokens. To get one:

1. Sign in at [cryptohopper.com](https://www.cryptohopper.com) and open the developer dashboard.
2. Create an OAuth application — you'll receive a `client_id` and `client_secret`.
3. Drive the OAuth consent flow (`/oauth-consent?app_id=<client_id>&redirect_uri=<your_uri>&state=<csrf>`) to receive a 40-character bearer token scoped to the permissions you requested.

Pass the token to `NewClient`. Optionally pass your OAuth `client_id` with `WithAppKey` — it's sent as the `x-api-app-key` header.

```go
client, err := cryptohopper.NewClient(
    os.Getenv("CRYPTOHOPPER_TOKEN"),
    cryptohopper.WithAppKey(os.Getenv("CRYPTOHOPPER_CLIENT_ID")),
)
```

## Options

| Option | Default | Description |
|---|---|---|
| `WithBaseURL(string)` | `https://api.cryptohopper.com/v1` | Override for staging |
| `WithHTTPClient(*http.Client)` | — | Inject a custom HTTP client |
| `WithTimeout(time.Duration)` | `30s` | Per-request timeout (ignored if `WithHTTPClient` set) |
| `WithUserAgent(string)` | — | Appended after `cryptohopper-go-sdk/<v>` |
| `WithAppKey(string)` | — | OAuth `client_id`, sent as `x-api-app-key` |
| `WithMaxRetries(int)` | `3` | Retries on HTTP 429. `0` disables auto-retry. |

## Errors

Every non-2xx response becomes a `*cryptohopper.Error`:

```go
var ce *cryptohopper.Error
if errors.As(err, &ce) {
    switch ce.Code {
    case "RATE_LIMITED":
        time.Sleep(ce.RetryAfter)
    case "UNAUTHORIZED":
        log.Fatal("token expired or invalid")
    case "FORBIDDEN":
        log.Printf("missing scope or IP mismatch (ip=%s)", ce.IPAddress)
    }
}
```

`ce.Code` values: `UNAUTHORIZED`, `FORBIDDEN`, `NOT_FOUND`, `RATE_LIMITED`, `VALIDATION_ERROR`, `DEVICE_UNAUTHORIZED`, `SERVER_ERROR`, `NETWORK_ERROR`, `TIMEOUT`, `UNKNOWN`. Unknown strings pass through verbatim.

## Rate limiting

The server enforces three buckets (`normal` 30/min, `order` 8/8s, `backtest` 1/2s). On 429 the SDK retries with exponential backoff up to `WithMaxRetries(n)` (default 3), honouring `Retry-After`.

## Development

```bash
go vet ./...
go test ./... -race -count=1
```

## Release

Push a `v*` git tag (plain `v` prefix required by Go modules). The workflow runs `go vet`, `go test -race`, and cuts a GitHub Release; Go proxies pick it up automatically from the tag.

## License

MIT — see [LICENSE](./LICENSE).
