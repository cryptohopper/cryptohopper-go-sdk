# Changelog

All notable changes to `github.com/cryptohopper/cryptohopper-go-sdk` are documented in this file.
The format is loosely based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## v0.2.0-alpha.1 — Unreleased

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
