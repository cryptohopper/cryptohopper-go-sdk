# Error Handling

Every non-2xx response and every transport failure returns a `*cryptohopper.Error`. The shape is the same idea as the Node/Python/Ruby/Rust/PHP/Dart SDKs but laid out as exported struct fields per Go idiom.

```go
type Error struct {
    Code       string         // shared SDK taxonomy
    Status     int            // HTTP status; 0 on transport failure
    Message    string         // server-provided human-readable
    ServerCode int            // numeric `code` from the envelope, 0 if absent
    IPAddress  string         // server-reported caller IP, empty if absent
    RetryAfter time.Duration  // parsed Retry-After (only on 429)
}
```

## Discriminating with `errors.As`

```go
import "errors"

if _, err := client.Hoppers.Get(ctx, "999999"); err != nil {
    var ce *cryptohopper.Error
    if errors.As(err, &ce) {
        log.Printf("code=%s status=%d server_code=%d ip=%s msg=%s",
            ce.Code, ce.Status, ce.ServerCode, ce.IPAddress, ce.Message)
    } else {
        log.Printf("non-SDK error: %v", err)
    }
}
```

`errors.As` follows wrapped errors, so this also works if you've wrapped the SDK error elsewhere with `fmt.Errorf("...: %w", err)`.

## Error code catalog

| `Code` | HTTP | When you'll see it | Recover by |
|---|---|---|---|
| `VALIDATION_ERROR` | 400, 422 | Missing or malformed parameter | Fix the request; the message says which parameter |
| `UNAUTHORIZED` | 401 | Token missing, wrong, or revoked | Re-auth |
| `DEVICE_UNAUTHORIZED` | 402 | Internal Cryptohopper device-auth flow rejected you | Shouldn't happen via the public API; contact support |
| `FORBIDDEN` | 403 | Scope missing, or IP not allowlisted | Check `ce.IPAddress`; add to allowlist or grant the scope |
| `NOT_FOUND` | 404 | Resource or endpoint doesn't exist | Check the ID; check you're using the latest SDK |
| `CONFLICT` | 409 | Resource is in a conflicting state | Cancel the existing job or wait |
| `RATE_LIMITED` | 429 | Bucket exhausted | The SDK auto-retries; see [Rate Limits](Rate-Limits.md) |
| `SERVER_ERROR` | 500–502, 504 | Cryptohopper's end | Retry with back-off; report if persistent |
| `SERVICE_UNAVAILABLE` | 503 | Planned maintenance or downstream outage | Respect `RetryAfter`; retry |
| `NETWORK_ERROR` | — | DNS failure, TCP reset, TLS handshake failure | Retry; check your network |
| `TIMEOUT` | — | Hit the client-side `WithTimeout` or `ctx.DeadlineExceeded` | Retry; bump timeout if the operation is legitimately slow |
| `UNKNOWN` | any | Anything else the SDK didn't recognise | Inspect `ce.Status` and `ce.Message` |

These strings are stable across SDK versions — compare with `==`, never substring-match.

## A switch statement with all the codes

```go
import "errors"

func categorize(err error) string {
    var ce *cryptohopper.Error
    if !errors.As(err, &ce) {
        return "non-sdk"
    }
    switch ce.Code {
    case "UNAUTHORIZED", "FORBIDDEN", "DEVICE_UNAUTHORIZED":
        return "auth"
    case "VALIDATION_ERROR":
        return "bad-request"
    case "NOT_FOUND":
        return "not-found"
    case "CONFLICT":
        return "conflict"
    case "RATE_LIMITED":
        return "throttled"
    case "SERVER_ERROR", "SERVICE_UNAVAILABLE":
        return "server"
    case "NETWORK_ERROR", "TIMEOUT":
        return "transient"
    default:
        return "unknown"
    }
}
```

`go vet` won't enforce exhaustiveness on a `string`-typed switch. Tools like [`exhaustive`](https://github.com/nishanths/exhaustive) won't help either since `Code` is just a string. If you want compile-time safety for a known subset, define your own typed enum and convert in one place:

```go
type Kind int

const (
    KindUnknown Kind = iota
    KindAuth
    KindRateLimited
    // ...
)

func categorize(err error) Kind { /* same switch */ }
```

## Context cancellation vs SDK errors

When `ctx.Done()` fires while a request is in flight, the SDK's transport returns the underlying context error wrapped as a `NETWORK_ERROR`-typed `*Error`. To distinguish a user-initiated cancellation from a real network failure:

```go
var ce *cryptohopper.Error
if errors.As(err, &ce) && ce.Code == "NETWORK_ERROR" {
    if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
        // user / parent context cancelled — drop quietly
        return nil
    }
    // genuine transport failure — log + retry
}
```

`errors.Is` checks the whole chain, so you'll catch the cancellation even though the visible top-level error is a `*cryptohopper.Error`.

## A robust retry wrapper

```go
func withRetry[T any](
    ctx context.Context,
    fn func(ctx context.Context) (T, error),
    maxAttempts int,
    baseDelay time.Duration,
) (T, error) {
    var zero T
    transient := map[string]bool{
        "SERVER_ERROR":        true,
        "SERVICE_UNAVAILABLE": true,
        "NETWORK_ERROR":       true,
        "TIMEOUT":              true,
    }
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        result, err := fn(ctx)
        if err == nil {
            return result, nil
        }
        var ce *cryptohopper.Error
        if !errors.As(err, &ce) || !transient[ce.Code] || attempt == maxAttempts {
            return zero, err
        }
        wait := baseDelay << (attempt - 1)
        if ce.RetryAfter > 0 {
            wait = ce.RetryAfter
        }
        select {
        case <-ctx.Done():
            return zero, ctx.Err()
        case <-time.After(wait):
        }
    }
    return zero, errors.New("unreachable")
}
```

Don't include `RATE_LIMITED` in `transient` — the SDK already retries 429s internally. Wrapping it here would multiply attempts unhelpfully.

## Logging

The `Error.Error()` method renders a compact one-line form including the code, status, message, and IP:

```
cryptohopper: [FORBIDDEN 403] IP not in allowlist (ip 203.0.113.5)
```

Suitable for direct logging via `log.Printf("%v", err)`. For structured logging, pull individual fields:

```go
import "log/slog"

var ce *cryptohopper.Error
if errors.As(err, &ce) {
    slog.Error("cryptohopper request failed",
        "code", ce.Code,
        "status", ce.Status,
        "server_code", ce.ServerCode,
        "ip", ce.IPAddress,
        "retry_after", ce.RetryAfter,
        "message", ce.Message,
    )
}
```
