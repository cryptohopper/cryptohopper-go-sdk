# Changelog

All notable changes to `github.com/cryptohopper/cryptohopper-go-sdk` are documented in this file.
The format is loosely based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## v0.4.0-alpha.2 — Unreleased

### Fixed
- **Critical: every authenticated request was rejected by the API gateway.** The transport sent `Authorization: Bearer <token>`, which the AWS API Gateway in front of `api.cryptohopper.com/v1/*` rejects (`405 Missing Authentication Token`). Cryptohopper's Public API v1 uses `access-token: <token>` — confirmed by the official [API documentation](https://www.cryptohopper.com/api-documentation/how-the-api-works) and the legacy iOS/Android SDKs. Switching to send `access-token` instead. The `Authorization` header is no longer set.

### Compatibility
No public-API change. Resource methods (`ch.User.Get(ctx)`, `ch.Hoppers.List(ctx, nil)`, etc.) keep their signatures. Only the wire-level header changes.

## v0.4.0-alpha.1 — 2026-04-24

Adds four more API domains: `Social`, `Tournaments`, `Webhooks`, `App`. Final A-wave — all 14 remaining public domains now covered.

### Added
- **`Social`** (27 methods) — profiles, feed, trends, search, notifications, conversations/messages, posts, comments, media, follows, likes/reposts, moderation.
- **`Tournaments`** (11 methods) — `List`, `Active`, `Get`, `Search`, `Trades`, `Stats`, `Activity`, `Leaderboard`, `TournamentLeaderboard`, `Join`, `Leave`.
- **`Webhooks`** (2 methods) — developer webhook registration (`/api/webhook_*`).
- **`App`** (2 methods) — mobile app store `Receipt` + `InAppPurchase`.

## v0.3.0-alpha.1 — 2026-04-24

Adds four more API domains: `AI`, `Platform`, `Chart`, `Subscription`.

### Added
- **`AI`** — `List`, `Get`, `AvailableModels`, `GetCredits`, `CreditInvoices`, `CreditTransactions`, `BuyCredits`, `LLMAnalyzeOptions`, `LLMAnalyze`, `LLMAnalyzeResults`, `LLMResults`.
- **`Platform`** — `LatestBlog`, `Documentation`, `PromoBar`, `SearchDocumentation`, `Countries`, `CountryAllowlist`, `IPCountry`, `Languages`, `BotTypes` (all public).
- **`Chart`** — `List`, `Get`, `Save`, `Delete`, `ShareSave`, `ShareGet`.
- **`Subscription`** — `Hopper`, `Get`, `Plans`, `Remap`, `Assign`, `GetCredits`, `OrderSub`, `StopSubscription`.

## v0.2.0-alpha.1 — 2026-04-24

Adds four more API domains: `Signals`, `Arbitrage`, `MarketMaker`, `Template`.

### Added
- **`Signals`** — `List`, `Performance`, `Stats`, `Distribution`, `ChartData` (signal-provider analytics; distinct from `Market.Signals` marketplace browse).
- **`Arbitrage`** — `ExchangeStart`, `ExchangeCancel`, `ExchangeResults`, `ExchangeHistory`, `ExchangeTotal`, `ExchangeResetTotal`, `MarketStart`, `MarketCancel`, `MarketResult`, `MarketHistory`, `Backlogs`, `Backlog`, `DeleteBacklog`.
- **`MarketMaker`** — `Get`, `Cancel`, `History`, `GetMarketTrend`, `SetMarketTrend`, `DeleteMarketTrend`, `Backlogs`, `Backlog`, `DeleteBacklog`.
- **`Template`** — `List`, `Get`, `Basic`, `Save`, `Update`, `Load`, `Delete`.

## v0.1.0-alpha.1 — 2026-04-24

Initial release. Covers six core API domains.

### Transport
- `Client` with functional options (`WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithUserAgent`, `WithAppKey`, `WithMaxRetries`).
- OAuth2 bearer auth via `Authorization: Bearer <token>`. Optional `x-api-app-key` via `WithAppKey`.
- `*Error` type with `Code`, `Status`, `Message`, `ServerCode`, `IPAddress`, `RetryAfter` — inspect via `errors.As(err, &ce)`.
- Automatic retry on HTTP 429 honouring `Retry-After` (default 3 retries, `WithMaxRetries(0)` to disable).
- Stdlib only — no third-party dependencies.

### Resources
- `User` — `Get`
- `Hoppers` — `List`, `Get`, `Create`, `Update`, `Delete`, `Positions`, `Position`, `Orders`, `Buy`, `Sell`, `ConfigGet`, `ConfigUpdate`, `ConfigPools`, `Panic`
- `Exchange` — `Ticker`, `Candles`, `Orderbook`, `Markets`, `Currencies`, `Exchanges`, `ForexRates`
- `Strategy` — `List`, `Get`, `Create`, `Update`, `Delete`
- `Backtest` — `Create`, `Get`, `List`, `Cancel`, `Restart`, `Limits`
- `Market` — `Signals`, `Signal`, `Items`, `Item`, `Homepage`
