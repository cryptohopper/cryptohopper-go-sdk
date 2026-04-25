# Getting Started

## Install

```bash
go get github.com/cryptohopper/cryptohopper-go-sdk@latest
```

Requires Go 1.22 or newer. The module has zero non-stdlib dependencies — `net/http`, `context`, `encoding/json`, and `time` carry the whole transport.

## First call

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cryptohopper "github.com/cryptohopper/cryptohopper-go-sdk"
)

func main() {
	client, err := cryptohopper.NewClient(os.Getenv("CRYPTOHOPPER_TOKEN"))
	if err != nil {
		log.Fatalf("client init: %v", err)
	}

	ctx := context.Background()
	me, err := client.User.Get(ctx)
	if err != nil {
		log.Fatalf("user.get: %v", err)
	}
	fmt.Printf("Logged in as: %s\n", me["email"])
}
```

Every method takes a `context.Context` as its first argument. Use `context.Background()` for top-level calls; pass a request-scoped context (e.g. from an HTTP handler) for downstream calls so cancellation propagates.

## Getting a token

Every request (except a handful of public endpoints like `/exchange/ticker`) needs an OAuth2 bearer token. Create one via **Developer → Create App** on [cryptohopper.com](https://www.cryptohopper.com) and complete the consent flow. The token is a 40-character opaque string.

For local dev:

```bash
export CRYPTOHOPPER_TOKEN=<your-token>
```

In production, load from your secret store of choice (AWS Secrets Manager, Vault, GCP Secret Manager, etc.) on startup.

## Idiomatic patterns

### Error handling with `errors.As`

The SDK returns a typed `*cryptohopper.Error` for every API failure. Use `errors.As` to discriminate:

```go
import "errors"

if _, err := client.Hoppers.Get(ctx, "999999"); err != nil {
	var ce *cryptohopper.Error
	if errors.As(err, &ce) {
		switch ce.Code {
		case "NOT_FOUND":
			// expected; ignore
		case "UNAUTHORIZED":
			// re-auth your token
			refresh()
		case "RATE_LIMITED":
			// SDK already retried; back off harder
			time.Sleep(ce.RetryAfter)
		default:
			log.Printf("cryptohopper: %s", err)
		}
	} else {
		// non-SDK error — log and surface
		return err
	}
}
```

### Cancellation through `context.Context`

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

_, err := client.Hoppers.List(ctx, nil)
if err != nil {
	if errors.Is(err, context.DeadlineExceeded) {
		log.Println("hopper list timed out — bumping ctx and retrying")
	}
	return err
}
```

### Functional options at construction

```go
client, err := cryptohopper.NewClient(
	os.Getenv("CRYPTOHOPPER_TOKEN"),
	cryptohopper.WithBaseURL("https://api.staging.cryptohopper.com/v1"),
	cryptohopper.WithTimeout(60*time.Second),
	cryptohopper.WithMaxRetries(5),
	cryptohopper.WithUserAgent("my-bot/1.2"),
	cryptohopper.WithAppKey(os.Getenv("CRYPTOHOPPER_APP_KEY")),
)
```

All options are independent. Pass any subset, in any order.

## Common pitfalls

**`cryptohopper: NewClient: apiKey must not be empty`** — you passed an empty string (often: `os.Getenv` for an unset variable). Check that the env var is actually exported in the process running your binary.

**`UNAUTHORIZED` on every call** — token is wrong, expired, or revoked. Visit the app page in the Cryptohopper dashboard to confirm.

**`FORBIDDEN` on endpoints that used to work** — IP allowlisting on the OAuth app blocked your current outbound IP. The error includes `IPAddress` so you can see what Cryptohopper saw:

```go
var ce *cryptohopper.Error
if errors.As(err, &ce) && ce.Code == "FORBIDDEN" {
	log.Printf("blocked from %s", ce.IPAddress)
}
```

**`x509: certificate signed by unknown authority`** — corporate proxy or self-signed root CA. Don't disable verification globally; build a custom `*http.Client` with the right CA bundle and inject it:

```go
caCertPool := x509.NewCertPool()
ca, _ := os.ReadFile("/path/to/corporate-ca.pem")
caCertPool.AppendCertsFromPEM(ca)

custom := &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: caCertPool},
	},
}

client, _ := cryptohopper.NewClient(token, cryptohopper.WithHTTPClient(custom))
```

**Goroutine leaks** — every method respects `ctx`. If you spawn many goroutines each making SDK calls, give them a context that's eventually cancelled (e.g. `errgroup.WithContext`) so in-flight requests don't outlive the surrounding work.

## Next steps

- [Authentication](Authentication.md) — bearer flow, app keys, IP whitelisting, custom HTTP clients
- [Error Handling](Error-Handling.md) — every error code, `errors.As` patterns, retry wrappers
- [Rate Limits](Rate-Limits.md) — auto-retry, customizing back-off, concurrency patterns
