package cryptohopper

import (
	"context"
	"fmt"
	"net/url"
)

// ArbitrageAPI is the resource namespace for exchange + market arbitrage ops.
// Two flavours: cross-exchange (Exchange*) and intra-exchange (Market*),
// plus a shared backlog surface.
type ArbitrageAPI struct {
	client *Client
}

// ─── Cross-exchange arbitrage ─────────────────────────────────────────────

// ExchangeStart begins a cross-exchange arbitrage run. Requires ``trade``.
func (a *ArbitrageAPI) ExchangeStart(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "POST", "/arbitrage/exchange", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ExchangeCancel cancels a cross-exchange arbitrage run. Requires ``trade``.
func (a *ArbitrageAPI) ExchangeCancel(ctx context.Context, body map[string]any) error {
	if body == nil {
		body = map[string]any{}
	}
	return a.client.request(ctx, "POST", "/arbitrage/cancel", nil, body, nil)
}

// ExchangeResults fetches exchange-arb results. Requires ``read``.
func (a *ArbitrageAPI) ExchangeResults(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := a.client.request(ctx, "GET", "/arbitrage/results", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ExchangeHistory returns historical exchange-arb runs. Requires ``read``.
func (a *ArbitrageAPI) ExchangeHistory(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := a.client.request(ctx, "GET", "/arbitrage/history", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ExchangeTotal returns running totals. Requires ``read``.
func (a *ArbitrageAPI) ExchangeTotal(ctx context.Context) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "GET", "/arbitrage/total", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ExchangeResetTotal resets the running totals. Requires ``manage``.
func (a *ArbitrageAPI) ExchangeResetTotal(ctx context.Context) error {
	return a.client.request(ctx, "POST", "/arbitrage/resettotal", nil, map[string]any{}, nil)
}

// ─── Intra-exchange market arbitrage ──────────────────────────────────────

// MarketStart begins an intra-exchange (e.g. triangular) arbitrage run.
// Requires ``trade``.
func (a *ArbitrageAPI) MarketStart(ctx context.Context, body map[string]any) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "POST", "/arbitrage/market", nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MarketCancel cancels a market-arb run. Requires ``trade``.
func (a *ArbitrageAPI) MarketCancel(ctx context.Context, body map[string]any) error {
	if body == nil {
		body = map[string]any{}
	}
	return a.client.request(ctx, "POST", "/arbitrage/market-cancel", nil, body, nil)
}

// MarketResult returns the result of a specific market-arb run. Requires ``read``.
func (a *ArbitrageAPI) MarketResult(ctx context.Context, extra url.Values) (map[string]any, error) {
	out := map[string]any{}
	if err := a.client.request(ctx, "GET", "/arbitrage/market-result", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MarketHistory returns historical market-arb runs. Requires ``read``.
func (a *ArbitrageAPI) MarketHistory(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := a.client.request(ctx, "GET", "/arbitrage/market-history", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Backlog (shared) ────────────────────────────────────────────────────

// Backlogs lists queued/pending backlog items. Requires ``read``.
func (a *ArbitrageAPI) Backlogs(ctx context.Context, extra url.Values) ([]map[string]any, error) {
	var out []map[string]any
	if err := a.client.request(ctx, "GET", "/arbitrage/get-backlogs", extra, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Backlog fetches a single backlog item. Requires ``read``.
func (a *ArbitrageAPI) Backlog(ctx context.Context, id any) (map[string]any, error) {
	q := url.Values{"backlog_id": {fmt.Sprint(id)}}
	out := map[string]any{}
	if err := a.client.request(ctx, "GET", "/arbitrage/get-backlog", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteBacklog deletes a backlog item. Requires ``manage``.
func (a *ArbitrageAPI) DeleteBacklog(ctx context.Context, id any) error {
	return a.client.request(ctx, "POST", "/arbitrage/delete-backlog", nil, map[string]any{"backlog_id": id}, nil)
}
