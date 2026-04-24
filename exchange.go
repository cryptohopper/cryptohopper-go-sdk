package cryptohopper

import (
	"context"
	"net/url"
	"strconv"
)

// ExchangeAPI is the resource namespace for public market data.
type ExchangeAPI struct {
	client *Client
}

// CandleOptions are the optional parameters for Candles.
type CandleOptions struct {
	From int64
	To   int64
}

// Ticker returns the current ticker for a market on an exchange.
func (e *ExchangeAPI) Ticker(ctx context.Context, exchange, market string) (map[string]any, error) {
	q := url.Values{"exchange": {exchange}, "market": {market}}
	out := map[string]any{}
	if err := e.client.request(ctx, "GET", "/exchange/ticker", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Candles returns OHLCV candles. Timeframe examples: "1m", "1h", "1d".
func (e *ExchangeAPI) Candles(ctx context.Context, exchange, market, timeframe string, opts *CandleOptions) ([]map[string]any, error) {
	q := url.Values{
		"exchange":  {exchange},
		"market":    {market},
		"timeframe": {timeframe},
	}
	if opts != nil {
		if opts.From > 0 {
			q.Set("from", strconv.FormatInt(opts.From, 10))
		}
		if opts.To > 0 {
			q.Set("to", strconv.FormatInt(opts.To, 10))
		}
	}
	var out []map[string]any
	if err := e.client.request(ctx, "GET", "/exchange/candle", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Orderbook returns the order book depth for a market.
func (e *ExchangeAPI) Orderbook(ctx context.Context, exchange, market string) (map[string]any, error) {
	q := url.Values{"exchange": {exchange}, "market": {market}}
	out := map[string]any{}
	if err := e.client.request(ctx, "GET", "/exchange/orderbook", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Markets lists the markets available on an exchange.
func (e *ExchangeAPI) Markets(ctx context.Context, exchange string) ([]map[string]any, error) {
	q := url.Values{"exchange": {exchange}}
	var out []map[string]any
	if err := e.client.request(ctx, "GET", "/exchange/markets", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Currencies lists the currencies available on an exchange.
func (e *ExchangeAPI) Currencies(ctx context.Context, exchange string) ([]map[string]any, error) {
	q := url.Values{"exchange": {exchange}}
	var out []map[string]any
	if err := e.client.request(ctx, "GET", "/exchange/currencies", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Exchanges lists all supported exchanges.
func (e *ExchangeAPI) Exchanges(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := e.client.request(ctx, "GET", "/exchange/exchanges", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ForexRates returns fiat forex rates used for conversion.
func (e *ExchangeAPI) ForexRates(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := e.client.request(ctx, "GET", "/exchange/forex-rates", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
