# Authentication

Every SDK request (except a handful of public endpoints) requires an OAuth2 bearer token:

```
Authorization: Bearer <40-char token>
```

## Obtaining a token

1. Log in to [cryptohopper.com](https://www.cryptohopper.com).
2. **Developer → Create App** — gives you a `client_id` + `client_secret`.
3. Complete the OAuth consent flow for your app, which returns a bearer token.

Options to automate step 3:

- **The official CLI**: `cryptohopper login` opens the consent page, runs a loopback listener, and persists the token to `~/.cryptohopper/config.json`. Read the token from there in your Go binary.
- **Your own code**: call the server's `/oauth2/authorize` + `/oauth2/token` endpoints directly. The CLI's implementation is short (~300 lines of TypeScript) and a reasonable reference.

## Client construction

```go
import cryptohopper "github.com/cryptohopper/cryptohopper-go-sdk"

client, err := cryptohopper.NewClient(
	os.Getenv("CRYPTOHOPPER_TOKEN"),
	cryptohopper.WithAppKey(os.Getenv("CRYPTOHOPPER_APP_KEY")),
	cryptohopper.WithBaseURL("https://api.cryptohopper.com/v1"),
	cryptohopper.WithTimeout(30*time.Second),
	cryptohopper.WithMaxRetries(3),
	cryptohopper.WithUserAgent("my-app/1.0"),
)
```

The `apiKey` argument is required; everything else is a `ClientOption` and optional.

### `WithAppKey`

Cryptohopper lets OAuth apps identify themselves on every request via the `x-api-app-key` header (value = your OAuth `client_id`). When set, the SDK adds the header automatically. Reasons to set it:

- Shows up in Cryptohopper's server-side telemetry, so you can attribute your own traffic.
- Drives per-app rate limits — if two apps share a token, they get independent quotas.
- Harmless to omit; the server accepts unattributed requests.

### `WithBaseURL`

Override for staging or a local dev server. The default is `https://api.cryptohopper.com/v1`. The trailing `/v1` is part of the base; resource paths are relative to it.

### `WithHTTPClient`

If you need custom transport behaviour — proxies, custom CA bundles, connection-pool tuning, OpenTelemetry instrumentation — supply your own `*http.Client`:

```go
custom := &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		Proxy: http.ProxyURL(corpProxyURL),
		TLSClientConfig: &tls.Config{RootCAs: corpRootCAs},
	},
}

client, err := cryptohopper.NewClient(
	token,
	cryptohopper.WithHTTPClient(custom),
)
```

The SDK does **not** take ownership of the client. You're responsible for any per-request cancellation, connection cleanup, etc.

To wire up tracing or per-call metrics, wrap `Transport` with a `RoundTripper` that records start/finish:

```go
custom := &http.Client{
	Timeout: 30 * time.Second,
	Transport: otelhttp.NewTransport(http.DefaultTransport),
}
```

### `WithTimeout`

Per-request timeout. Defaults to 30 seconds. The 429-retry path may stack additional time on top of this — set it conservatively if your `WithMaxRetries` is high.

### `WithMaxRetries`

Number of automatic retries on HTTP 429. Default 3. Set to 0 to disable. See [Rate Limits](Rate-Limits.md) for details.

### `WithUserAgent`

Appended after the SDK's own User-Agent (`cryptohopper-sdk-go/<version>`). Set this to identify your client to support if you ever need to debug something with Cryptohopper.

## IP allowlisting

If your Cryptohopper app has IP allowlisting enabled, requests from unlisted IPs return `403 FORBIDDEN`. The SDK surfaces this as `*cryptohopper.Error` with `Code == "FORBIDDEN"` and a populated `IPAddress` field showing the IP Cryptohopper saw:

```go
import "errors"

var ce *cryptohopper.Error
if errors.As(err, &ce) && ce.Code == "FORBIDDEN" {
	log.Printf("blocked: %s", ce.IPAddress)
}
```

For CI where the runner IP isn't stable, either disable IP allowlisting for that app or route outbound traffic through a stable IP (NAT gateway, VPN, dedicated proxy).

## Rotating tokens

Cryptohopper bearer tokens are long-lived but can be revoked:

- Manually from the dashboard.
- When the user revokes consent.

The SDK surfaces revocation as `UNAUTHORIZED` on the next call. There is no automatic refresh-token handling in the SDK today — if your app uses refresh tokens, handle the `UNAUTHORIZED` branch by exchanging your refresh token for a new access token, then constructing a fresh client:

```go
func withAutoRefresh[T any](
	ctx context.Context,
	mu *sync.Mutex,
	clientPtr **cryptohopper.Client,
	fn func(*cryptohopper.Client) (T, error),
) (T, error) {
	mu.Lock()
	c := *clientPtr
	mu.Unlock()

	result, err := fn(c)
	if err == nil {
		return result, nil
	}

	var ce *cryptohopper.Error
	if !errors.As(err, &ce) || ce.Code != "UNAUTHORIZED" {
		return result, err
	}

	mu.Lock()
	defer mu.Unlock()
	newToken, err := refreshTokenFromYourStore()
	if err != nil {
		return result, err
	}
	newClient, err := cryptohopper.NewClient(newToken /* + same options */)
	if err != nil {
		return result, err
	}
	*clientPtr = newClient
	return fn(newClient)
}
```

The client's API key is intentionally immutable. Construct a fresh client for token rotation; the cost is small and it sidesteps races where one in-flight request uses an old token while another uses the new.

## Concurrency

`*Client` is safe for concurrent use across goroutines. The underlying `*http.Client` is shared and pooled. You don't need a client-per-goroutine; one client serving an `errgroup` or worker pool is fine.

```go
import "golang.org/x/sync/errgroup"

g, gctx := errgroup.WithContext(ctx)
for _, id := range hopperIDs {
	id := id
	g.Go(func() error {
		_, err := client.Hoppers.Get(gctx, id)
		return err
	})
}
if err := g.Wait(); err != nil {
	return err
}
```

See [Rate Limits](Rate-Limits.md) for guidance on capping concurrency at the API quota.

## Public-only access (no token)

A handful of endpoints accept anonymous calls:

- `/market/*` — marketplace browse
- `/platform/*` — i18n, country list, blog feed
- `/exchange/ticker`, `/exchange/candle`, `/exchange/orderbook`, `/exchange/markets`, `/exchange/exchanges`, `/exchange/forex-rates` — public market data

The SDK still requires `apiKey` at construction; pass any non-empty placeholder if you only intend to hit public endpoints. The server ignores the bearer header on whitelisted routes.

```go
client, _ := cryptohopper.NewClient("anonymous")
btc, _ := client.Exchange.Ticker(ctx, &cryptohopper.TickerParams{
	Exchange: "binance",
	Market:   "BTC/USDT",
})
```
