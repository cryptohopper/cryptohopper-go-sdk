# Rate Limits

Cryptohopper applies per-bucket rate limits on the server. When you hit one, you get a `429` with a `Retry-After` header. The SDK handles this for you.

## The default behaviour

On every `429`, the SDK:

1. Parses `Retry-After` into a `time.Duration` (handles both seconds-as-integer and HTTP-date forms).
2. Sleeps that long (falling back to exponential back-off if the header is missing).
3. Retries the request.
4. Repeats up to `WithMaxRetries(n)` (default 3).

If retries exhaust, the call returns a `*cryptohopper.Error` with `Code == "RATE_LIMITED"` and `RetryAfter` set to the last seen retry hint.

## Configuring it

```go
client, _ := cryptohopper.NewClient(
    token,
    cryptohopper.WithMaxRetries(10),
    cryptohopper.WithTimeout(60*time.Second),
)
```

To **disable** retries entirely (e.g. you want to do your own back-off):

```go
client, _ := cryptohopper.NewClient(
    token,
    cryptohopper.WithMaxRetries(0),
)
```

With `WithMaxRetries(0)` a 429 surfaces immediately as `RATE_LIMITED`. Inspect `ce.RetryAfter` (a `time.Duration`) and schedule the retry on your own timeline.

## Buckets

Cryptohopper has three named buckets:

| Bucket | Scope | Example endpoints |
|---|---|---|
| `normal` | Most reads + writes | `/user/get`, `/hopper/list`, `/hopper/update`, `/exchange/ticker` |
| `order` | Anything that places or modifies orders | `/hopper/buy`, `/hopper/sell`, `/hopper/panic` |
| `backtest` | The (expensive) backtest subsystem | `/backtest/new`, `/backtest/get` |

The SDK doesn't know which bucket a call hits — it only sees the 429. You don't need to either; the server tells you when you're limited.

## Backfill jobs (own back-off)

If you're ingesting historical data and need to fetch many pages, take ownership of the back-off:

```go
client, _ := cryptohopper.NewClient(token, cryptohopper.WithMaxRetries(0))

for _, hopperID := range allHopperIDs {
    for {
        orders, err := client.Hoppers.Orders(ctx, hopperID, nil)
        if err == nil {
            process(orders)
            break
        }
        var ce *cryptohopper.Error
        if errors.As(err, &ce) && ce.Code == "RATE_LIMITED" {
            wait := ce.RetryAfter
            if wait == 0 {
                wait = time.Second
            }
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(wait):
            }
            continue
        }
        return err
    }
}
return nil
```

This pattern lets a long-running job honour rate limits without stalling other work, because you decide the pacing.

## Concurrency caps with `errgroup` + semaphore

```go
import (
    "golang.org/x/sync/errgroup"
    "golang.org/x/sync/semaphore"
)

const maxConcurrent = 4

sem := semaphore.NewWeighted(maxConcurrent)
g, gctx := errgroup.WithContext(ctx)

for _, id := range hopperIDs {
    id := id
    if err := sem.Acquire(gctx, 1); err != nil {
        return err
    }
    g.Go(func() error {
        defer sem.Release(1)
        _, err := client.Hoppers.Get(gctx, id)
        return err
    })
}
if err := g.Wait(); err != nil {
    return err
}
```

Empirically, **4–8 concurrent workers** is comfortable for most accounts. Higher is feasible with `WithAppKey` set (which gives your OAuth app its own quota) but plan to back off explicitly.

## What the SDK does NOT do

- **No global semaphore.** If you spawn 100 goroutines each calling the SDK and the server rate-limits them, every goroutine's retry is independent — you might get 100 simultaneous sleeps. Cap concurrency yourself.
- **No adaptive slow-down.** After a 429, the SDK waits and retries that one call. It doesn't throttle future calls. If you see frequent 429s, lower your concurrency or add explicit pacing.
- **No client-side bucket tracking.** The server is the source of truth.

## Diagnosing "always rate-limited"

If every request returns `RATE_LIMITED` even at low volume:

1. Check that your app hasn't been flagged for abuse in the Cryptohopper dashboard.
2. Confirm you haven't accidentally created a loop that retries on non-429 errors too — `errors.As` discriminates correctly only if `Code == "RATE_LIMITED"`.
3. Inspect `ce.ServerCode` — Cryptohopper sometimes includes a numeric detail there that clarifies which bucket you've tripped.
4. Confirm you're not sharing one token across many machines (one quota, divided across all of them). If you have multiple environments, give each a distinct token.

## Interaction with `context.Context` deadlines

If `ctx` has a deadline and the SDK is mid-back-off when it expires, the retry sleep is interrupted via the SDK's internal `ctx.Done()` watch and the call returns a `NETWORK_ERROR` wrapping `context.DeadlineExceeded`. Use `errors.Is(err, context.DeadlineExceeded)` to detect this distinct from "server is genuinely slow".

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if _, err := client.Hoppers.List(ctx, nil); err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("ran out of time during back-off — try a longer ctx deadline")
    }
}
```
